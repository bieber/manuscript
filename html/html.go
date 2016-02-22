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
	"github.com/dustin/go-humanize"
	"io"
)

// Renderer provides a Render method to render the given document to
// an HTML file.
type Renderer struct {
	styleSheet    string
	authorInfo    bool
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
	authorInfo := false

	for k, v := range options {
		switch k {
		case "styleSheet":
			styleSheet = v
		case "authorInfo":
			authorInfo = argIsTrue(v)
		default:
			return nil, fmt.Errorf("Invalid HTML option %s", k)
		}
	}

	return &Renderer{
		styleSheet: styleSheet,
		authorInfo: authorInfo,
		document:   document,
	}, nil
}

// Render writes the requested document out to the specified io.Writer
// as an HTML file.
func (r *Renderer) Render(fout io.Writer) error {
	encoder := xml.NewEncoder(selfClosingRemover{fout})

	bodyContents := []interface{}{}
	bodyContents = append(bodyContents, r.renderFrontMatter())

	encoder.Indent("", "\t")
	return encoder.Encode(
		document{
			Head: r.renderHead(),
			Body: body{
				Content: div{
					Class:    "container",
					Children: bodyContents,
				},
			},
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

func (r *Renderer) renderFrontMatter() div {
	document := r.document

	contents := []interface{}{}

	if r.authorInfo {
		authorContents := []interface{}{}
		if document.Author.Name != "" {
			authorContents = append(
				authorContents,
				span{Text: document.Author.Name},
				br{},
			)
		}
		if len(document.Author.Address) != 0 {
			for _, l := range document.Author.Address {
				authorContents = append(
					authorContents,
					span{Text: l},
					br{},
				)
			}
		}
		if document.Author.PhoneNumber != "" {
			authorContents = append(
				authorContents,
				span{Text: document.Author.PhoneNumber},
				br{},
			)
		}
		if document.Author.EmailAddress != "" {
			authorContents = append(
				authorContents,
				span{Text: document.Author.EmailAddress},
				br{},
			)
		}
		if len(document.Author.ProfessionalOrgs) != 0 {
			for _, l := range document.Author.ProfessionalOrgs {
				authorContents = append(
					authorContents,
					span{Text: l},
					br{},
				)
			}
		}

		contents = append(
			contents,
			div{Class: "author_info", Children: authorContents},
		)
	}

	contents = append(contents, h1{Title: document.Title})

	authorText := "by " + document.Author.Byline
	if r.document.Type == parser.Novel {
		authorText = "a novel " + authorText
	}
	contents = append(contents, p{Class: "byline", Text: authorText})

	wordText := "about " + humanize.Comma(document.WordCount()) + " words"
	contents = append(contents, p{Class: "word_count", Text: wordText})

	return div{
		Class:    "front_matter",
		Children: contents,
	}
}
