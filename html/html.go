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
	"encoding/xml"
	"fmt"
	"github.com/bieber/manuscript/parser"
	"github.com/bieber/manuscript/renderers"
	"io"
)

// Renderer provides a Render method to render the given document to
// an HTML file.
type Renderer struct {
	styleSheet    string
	partNumber    int
	chapterNumber int
	lastElement   parser.DocumentElement
	document      parser.Document
}

// New constructs a new Renderer for the given document and
// command-line arguments.
func New(
	document parser.Document,
	options map[string]string,
) (renderers.Renderer, error) {
	styleSheet := ""

	for k, v := range options {
		switch k {
		case "styleSheet":
			styleSheet = v
		default:
			return nil, fmt.Errorf("Invalid HTMLL option %s", k)
		}
	}

	return &Renderer{
		styleSheet: styleSheet,
		document:   document,
	}, nil
}

// Render writes the requested document out to the specified io.Writer
// as an HTML file.
func (r *Renderer) Render(fout io.Writer) error {
	encoder := xml.NewEncoder(fout)

	encoder.Indent("", "\t")
	return encoder.Encode(
		document{
			Head: r.renderHead(),
		},
	)
}

func (r *Renderer) renderHead() header {
	var styleSheet *link
	if r.styleSheet != "" {
		styleSheet = &link{
			Rel:  "stylesheet",
			Type: "text/css",
			HREF: r.styleSheet,
		}
	}
	return header{
		Title:      r.document.Title,
		StyleSheet: styleSheet,
	}
}
