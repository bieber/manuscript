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

package bbcode

import (
	"errors"
	"fmt"
	"github.com/StefanSchroeder/Golang-Roman"
	"github.com/bieber/manuscript/parser"
	"github.com/bieber/manuscript/renderers"
	"io"
	"bytes"
)

// Renderer provides a Render method to render the given document to
// bbcode text.
type Renderer struct {
	document parser.Document
	buffer   bytes.Buffer
}

// New constructs a new Renderer for the given document and
// command-line arguments.
func New(
	document parser.Document,
	options map[string]string,
) (renderers.Renderer, error) {
	return &Renderer{document: document}, nil
}

// Render writes the requested document out to the specified io.Writer
// as bbcode text.
func (r *Renderer) Render(fout io.Writer) error {
	for _, p := range r.document.Parts {
		err := r.renderPart(p)
		if err != nil {
			return err
		}
	}

	_, err := r.buffer.WriteTo(fout)
	return err
}

func (r *Renderer) renderPart(part parser.Part) error {
	if !part.Anonymous {
		text := "Part " + roman.Roman(part.Number)
		if part.Title != "" {
			text += ": " + part.Title
		}

		_, err := r.buffer.WriteString("[b]" + text + "[/b]\n\n")
		if err != nil {
			return err
		}
	}

	for _, c := range part.Chapters {
		err := r.renderChapter(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Renderer) renderChapter(chapter parser.Chapter) error {
	if !chapter.Anonymous {
		text := ""
		if chapter.Prologue {
			text = "Prologue"
			if chapter.Title != "" {
				text += ": " + chapter.Title
			}
		} else {
			text = fmt.Sprintf("Chapter %d", chapter.Number)
			if chapter.Title != "" {
				text += ": " + chapter.Title
			}
		}

		_, err := r.buffer.WriteString("[b]" + text + "[/b]\n\n")
		if err != nil {
			return err
		}
	}

	for i, s := range chapter.Scenes {
		err := r.renderScene(s)
		if err != nil {
			return err
		}

		if i != len(chapter.Scenes)-1 {
			_, err := r.buffer.WriteString("------\n\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Renderer) renderScene(scene parser.Scene) error {
	for _, p := range scene.Paragraphs {
		err := r.renderParagraph(p)
		if err != nil {
			return err
		}

		_, err = r.buffer.WriteString("\n\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) renderParagraph(paragraph parser.Paragraph) error {
	for _, e := range paragraph.Text {
		err := r.renderElement(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) renderElement(element parser.DocumentElement) error {
	var err error
	switch e := element.(type) {
	case parser.PlainText:
		_, err = r.buffer.WriteString(string(e))
	case parser.ItalicText:
		_, err = r.buffer.WriteString("[i]" + string(e) + "[/i]")
	case parser.BoldText:
		_, err = r.buffer.WriteString("[b]" + string(e) + "[/b]")
	case parser.BoldItalicText:
		_, err = r.buffer.WriteString("[b][i]" + string(e) + "[/i][/b]")
	default:
		panic(
			errors.New(
				"bbcode: Unexpected document element passed to renderElement",
			),
		)
	}

	return err
}
