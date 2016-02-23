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

package html

import (
	"bytes"
	"io"
	"strings"
)

func argIsTrue(arg string) bool {
	arg = strings.ToLower(arg)
	return arg == "t" || arg == "true" || arg == "yes" || arg == "y"
}

type selfClosingRemover struct {
	dest io.Writer
}

// This is a pretty fragile implementation that will definitely break
// if Write ever comes in batches and one of them straddles a
// self-closing close tag.  I'm hoping that won't happen (it looks
// like the XML encoder just dumps the whole file at once), but if it
// does then I'll need to come back here and do some reasonable
// buffering.
func (s selfClosingRemover) Write(p []byte) (n int, err error) {
	n = len(p)
	toRemove := []string{
		"br",
		"link",
	}

	for _, tag := range toRemove {
		p = bytes.Replace(p, []byte("</"+tag+">"), []byte{}, -1)
	}

	_, err = s.dest.Write(p)
	return
}
