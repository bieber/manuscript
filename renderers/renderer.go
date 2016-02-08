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

package renderers

import (
	"github.com/bieber/manuscript/parser"
	"io"
)

// RendererConstructor is a function that creates a new Renderer given
// a document and a set of options as string key/value pairs.
type RendererConstructor func(parser.Document, map[string]string) Renderer

// Renderer defines an object capable of rendering a document.
// to the given output.
type Renderer interface {
	Render(io.Writer) error
}
