/* Copyright (c) 2016 Robert Bieber
 *
 * This file is part of manuscript.
 *
 * manuscript is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package util

import (
	"fmt"
	"github.com/StefanSchroeder/Golang-Roman"
)

// PartLabel assembles a label for a document part.
func PartLabel(number int, title string) string {
	text := "Part " + roman.Roman(number)
	if title != "" {
		text += ": " + title
	}
	return text
}

// PrologueLabel assembles a label for a prologue.
func PrologueLabel(title string) string {
	text := "Prologue"
	if title != "" {
		text += ": " + title
	}
	return text
}

// ChapterLabel assembles a label for a chapter.
func ChapterLabel(number int, title string) string {
	text := fmt.Sprintf("Chapter %d", number)
	if title != "" {
		text += ": " + title
	}
	return text
}
