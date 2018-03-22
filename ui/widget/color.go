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

package widget

import (
	"fmt"
	"hash/fnv"
	"sort"

	"github.com/gdamore/tcell"
)

var colorNames []string

// init initializes the colorNames array.
func init() {
	colorNames = make([]string, len(tcell.ColorNames))
	i := 0
	for name, _ := range tcell.ColorNames {
		colorNames[i] = name
		i++
	}
	// In order to have consistent coloring between restarts, we need to sort the array.
	sort.Sort(sort.StringSlice(colorNames))
}

// GetHashColorName gets a color name for the given string based on its FNV-1 hash.
//
// The array of possible color names are the alphabetically ordered color
// names specified in tcell.ColorNames.
//
// The algorithm to get the color is as follows:
//  colorNames[ FNV1(string) % len(colorNames) ]
//
// With the exception of the three special cases:
//  --> = green
//  <-- = red
//  --- = yellow
func GetHashColorName(s string) string {
	switch s {
	case "-->":
		return "green"
	case "<--":
		return "red"
	case "---":
		return "yellow"
	default:
		h := fnv.New32a()
		h.Write([]byte(s))
		return colorNames[int(h.Sum32())%len(colorNames)]
	}
}

// GetHashColor gets the tcell Color value for the given string.
//
// GetHashColor calls GetHashColorName() and gets the Color value from the tcell.ColorNames map.
func GetHashColor(s string) tcell.Color {
	return tcell.ColorNames[GetHashColorName(s)]
}

// AddHashColor adds tview color tags to the given string.
// The color added is the color returned by GetHashColorName().
func AddHashColor(s string) string {
	return fmt.Sprintf("[%s]%s[white]", GetHashColorName(s), s)
}
