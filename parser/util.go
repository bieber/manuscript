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

package parser

import (
	"math"
	"strings"
)

// WordCount returns an approximate word count for the document,
// rounded to the nearest 100 words for stories < 15,000 words, and to
// the nearest 500 for anything longer.
func (d Document) WordCount() int64 {
	count := 0
	for _, e := range d.Text {
		switch e := e.(type) {
		case PlainText:
			count += len(strings.Split(string(e), " "))
		case ItalicText:
			count += len(strings.Split(string(e), " "))
		case BoldText:
			count += len(strings.Split(string(e), " "))
		case BoldItalicText:
			count += len(strings.Split(string(e), " "))
		}
	}

	granularity := 100.0
	if count > 15000 {
		granularity = 500.0
	}
	return int64(granularity * math.Floor((float64(count)/granularity)+0.5))
}
