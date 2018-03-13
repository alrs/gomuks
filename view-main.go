// gomuks - A terminal Matrix client written in Go.
// Copyright (C) 2018 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type RoomView struct {
	*tview.Grid

	sender, message *tview.TextView
}

func NewRoomView() *RoomView {
	view := &RoomView{
		tview.NewGrid(),
		tview.NewTextView(),
		tview.NewTextView(),
	}
	view.SetColumns(30, 0).SetRows(0)

	view.sender.SetTextAlign(tview.AlignRight)
	view.sender.SetScrollable(true)
	view.message.SetScrollable(true)

	view.AddItem(view.sender, 0, 0, 1, 1, 1, 1, false)
	view.AddItem(view.message, 0, 1, 1, 1, 1, 1, false)

	return view
}

func (ui *GomuksUI) MakeMainUI() tview.Primitive {
	ui.mainView = tview.NewGrid().SetColumns(40, 0).SetRows(0, 2)

	ui.mainViewRoomList = tview.NewList().ShowSecondaryText(false)
	ui.mainViewRoomList.SetBorderPadding(1, 1, 1, 1)
	ui.mainView.AddItem(ui.mainViewRoomList, 0, 0, 2, 1, 2, 1, false)

	ui.mainViewRoomView = tview.NewPages()
	ui.mainViewRoomView.SetChangedFunc(ui.Render)
	ui.mainView.AddItem(ui.mainViewRoomView, 0, 1, 1, 1, 1, 1, false)

	ui.mainViewInput = tview.NewInputField()
	ui.mainViewInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			room, text := ui.currentRoom(), ui.mainViewInput.GetText()
			if len(text) == 0 {
				return
			} else if text[0] == '/' {
				args := strings.SplitN(text, " ", 2)
				command := strings.ToLower(args[0])
				args = args[1:]
				ui.HandleCommand(room, command, args)
			} else {
				ui.matrix.SendMessage(room, text)
			}
			ui.mainViewInput.SetText("")
		}
	})
	ui.mainView.AddItem(ui.mainViewInput, 1, 1, 1, 1, 1, 1, true)

	ui.debug.Print(ui.mainViewInput.SetInputCapture(ui.MainUIKeyHandler))

	ui.mainViewRooms = make(map[string]*RoomView)

	return ui.mainView
}

func (ui *GomuksUI) HandleCommand(room, command string, args []string) {
	switch command {
	case "quit":
		ui.matrix.Stop()
		ui.app.Stop()
	case "part":
	case "leave":
		ui.matrix.client.LeaveRoom(room)
	case "join":
		if len(args) == 0 {
			ui.Append(room, "*", "Usage: /join <room>")
		}
		mxid := args[0]
		server := mxid[strings.Index(mxid, ":")+1:]
		ui.matrix.client.JoinRoom(mxid, server, nil)
	}
}

func (ui *GomuksUI) MainUIKeyHandler(key *tcell.EventKey) *tcell.EventKey {
	ui.debug.Print(key)
	if key.Modifiers() == tcell.ModCtrl {
		if key.Key() == tcell.KeyDown {
			ui.SwitchRoom(ui.currentRoomIndex + 1)
			ui.mainViewRoomList.SetCurrentItem(ui.currentRoomIndex)
		} else if key.Key() == tcell.KeyUp {
			ui.SwitchRoom(ui.currentRoomIndex - 1)
			ui.mainViewRoomList.SetCurrentItem(ui.currentRoomIndex)
		}
	} else if key.Key() == tcell.KeyPgUp || key.Key() == tcell.KeyPgDn {
		ui.mainViewRooms[ui.currentRoom()].sender.InputHandler()(key, nil)
		ui.mainViewRooms[ui.currentRoom()].message.InputHandler()(key, nil)
	} else {
		return key
	}
	return nil
}

func (ui *GomuksUI) SetRoomList(rooms []string) {
	ui.roomList = rooms
	ui.mainViewRoomList.Clear()
	for index, room := range rooms {
		localRoomIndex := index
		ui.mainViewRoomList.AddItem(room, "", 0, func() {
			ui.SwitchRoom(localRoomIndex)
		})
		if !ui.mainViewRoomView.HasPage(room) {
			roomView := NewRoomView()
			ui.mainViewRooms[room] = roomView
			ui.mainViewRoomView.AddPage(room, roomView, true, false)
		}
	}
	ui.SwitchRoom(0)
}

func (ui *GomuksUI) currentRoom() string {
	if len(ui.roomList) == 0 {
		return ""
	}
	return ui.roomList[ui.currentRoomIndex]
}

func (ui *GomuksUI) SwitchRoom(roomIndex int) {
	if roomIndex < 0 {
		roomIndex = len(ui.roomList) - 1
	}
	ui.currentRoomIndex = roomIndex % len(ui.roomList)
	ui.mainViewRoomView.SwitchToPage(ui.roomList[ui.currentRoomIndex])
}

func (ui *GomuksUI) Append(room, sender, message string) {
	roomView, ok := ui.mainViewRooms[room]
	if ok {
		fmt.Fprintf(roomView.sender, sender)
		fmt.Fprintf(roomView.message, sender)
		ui.Render()
	}
}
