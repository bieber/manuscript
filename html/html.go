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
	"github.com/StefanSchroeder/Golang-Roman"
	"github.com/bieber/manuscript/parser"
	"github.com/bieber/manuscript/renderers"
	"github.com/dustin/go-humanize"
	"io"
)

// Renderer provides a Render method to render the given document to
// an HTML file.
type Renderer struct {
	styleSheet     string
	authorInfo     bool
	includeTOC     bool
	prologues      []string
	parts          []string
	chapters       []string
	partNumber     int
	chapterNumber  int
	prologueNumber int
	textPosition   int
	document       parser.Document
}

// New constructs a new Renderer for the given document and
// command-line arguments.
func New(
	document parser.Document,
	options map[string]string,
) (renderers.Renderer, error) {
	renderer := Renderer{
		document: document,
	}

	for k, v := range options {
		switch k {
		case "styleSheet":
			renderer.styleSheet = v
		case "authorInfo":
			renderer.authorInfo = argIsTrue(v)
		case "includeTOC":
			renderer.includeTOC = argIsTrue(v)
		default:
			return nil, fmt.Errorf("Invalid HTML option %s", k)
		}
	}

	for _, element := range document.Text {
		switch e := element.(type) {
		case parser.PrologueBreak:
			renderer.prologues = append(renderer.prologues, string(e))

		case parser.PartBreak:
			renderer.parts = append(renderer.parts, string(e))

		case parser.ChapterBreak:
			renderer.chapters = append(renderer.chapters, string(e))
		}
	}

	return &renderer, nil
}

// Render writes the requested document out to the specified io.Writer
// as an HTML file.
func (r *Renderer) Render(fout io.Writer) error {
	encoder := xml.NewEncoder(selfClosingRemover{fout})

	bodyContents := []interface{}{}
	bodyContents = append(bodyContents, r.renderFrontMatter())

	if r.includeTOC {
		bodyContents = append(bodyContents, r.renderTOC())
	}

	for r.textPosition < len(r.document.Text) {
		bodyContents = append(bodyContents, r.renderPart())
	}

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

func (r *Renderer) renderTOC() div {
	partNumber, chapterNumber, prologueNumber := 0, 0, 0
	part := li{}
	listChildren := []interface{}{}

	addNode := func(node interface{}) {
		// Indicates that we actually have one or more parts
		if len(part.Children) != 0 {
			if len(part.Children) < 2 {
				// There is no existing list of sub-items in this part yet
				part.Children = append(
					part.Children,
					ol{
						Class:    "toc_inner",
						Children: []interface{}{node},
					},
				)
			} else {
				// There is, so just add to it
				existingChildren := part.Children[1].(ol).Children
				part.Children[1] = ol{
					Class:    "toc_inner",
					Children: append(existingChildren, node),
				}
			}
		} else {
			listChildren = append(listChildren, node)
		}
	}

	addPart := func(text string) {
		if len(part.Children) != 0 {
			listChildren = append(listChildren, part)
			part = li{}
		}
		part.Children = append(
			part.Children,
			a{
				HREF: fmt.Sprintf("#part_%d", partNumber),
				Text: text,
			},
		)
	}

	for _, element := range r.document.Text {
		switch e := element.(type) {

		case parser.PartBreak:
			partNumber++
			chapterNumber = 0

			text := "Part " + roman.Roman(partNumber)
			if e != "" {
				text += ": " + string(e)
			}
			addPart(text)

		case parser.PrologueBreak:
			prologueNumber++

			text := "Prologue"
			if e != "" {
				text += ": " + string(e)
			}

			addNode(
				li{
					Children: []interface{}{
						a{
							HREF: fmt.Sprintf("#prologue_%d", prologueNumber),
							Text: text,
						},
					},
				},
			)

		case parser.ChapterBreak:
			chapterNumber++

			text := fmt.Sprintf("Chapter %d", chapterNumber)
			if e != "" {
				text += ": " + string(e)
			}

			addNode(
				li{
					Children: []interface{}{
						a{
							HREF: fmt.Sprintf(
								"#chapter_%d_%d",
								partNumber,
								chapterNumber,
							),
							Text: text,
						},
					},
				},
			)
		}
	}

	if len(part.Children) != 0 {
		listChildren = append(listChildren, part)
	}

	return div{
		Class:    "table_of_contents",
		Children: []interface{}{ol{Class: "toc_outer", Children: listChildren}},
	}
}

func (r *Renderer) renderPart() div {
	// We may not always have an explicit part marker (short stories,
	// novels without parts, or just prologue text before the first
	// part of a novel).  For these cases, we'll still render using
	// the regular procedure for a part but just leave out the header
	// and call it an "anonymous part."
	class := "anonymous_part"
	children := []interface{}{}

	if e, ok := r.peekElement().(parser.PartBreak); ok {
		r.partNumber++
		r.chapterNumber = 0

		class = "part"
		text := "Part " + roman.Roman(r.partNumber)
		if e != "" {
			text += ": " + string(e)
		}

		children = append(
			children,
			h2{
				Children: []interface{}{
					a{
						Name: fmt.Sprintf("part_%d", r.partNumber),
						Text: text,
					},
				},
			},
		)
		r.nextElement()
	}

	for r.hasNextElement() {
		if _, ok := r.peekElement().(parser.PartBreak); ok {
			break
		}

		children = append(
			children,
			r.renderChapter(),
		)
	}

	return div{
		Class:    class,
		Children: children,
	}

}

func (r *Renderer) renderChapter() div {
	// We may not always have an explicit chapter or prologue marker
	// (short stories, or just loose text at the beginning of a story
	// or part).  In that case we still render using the normal method
	// for a chapter, but call it an "anonymous chapter."
	class := "anonymous_chapter"
	children := []interface{}{}

	switch e := r.peekElement().(type) {
	case parser.PrologueBreak:
		r.prologueNumber++
		class = "chapter prologue"

		text := "Prologue"
		if e != "" {
			text += ": " + string(e)
		}

		children = append(
			children,
			h3{
				Children: []interface{}{
					a{
						Name: fmt.Sprintf("prologue_%d", r.prologueNumber),
						Text: text,
					},
				},
			},
		)

		r.nextElement()

	case parser.ChapterBreak:
		r.chapterNumber++
		class = "chapter"

		text := fmt.Sprintf("Chapter %d", r.chapterNumber)
		if e != "" {
			text = ": " + string(e)
		}

		children = append(
			children,
			h3{
				Children: []interface{}{
					a{
						Name: fmt.Sprintf(
							"chapter_%d_%d",
							r.partNumber,
							r.chapterNumber,
						),
						Text: text,
					},
				},
			},
		)

		r.nextElement()
	}

outer:
	for r.hasNextElement() {
		switch r.peekElement().(type) {
		case parser.PartBreak:
			break outer
		case parser.PrologueBreak:
			break outer
		case parser.ChapterBreak:
			break outer
		}

		r.nextElement()
	}

	return div{
		Class:    class,
		Children: children,
	}
}

func (r *Renderer) hasNextElement() bool {
	return r.textPosition < len(r.document.Text)
}

func (r *Renderer) peekElement() parser.DocumentElement {
	return r.document.Text[r.textPosition]
}

func (r *Renderer) nextElement() parser.DocumentElement {
	tp := r.textPosition
	r.textPosition++
	return r.document.Text[tp]
}

func (r *Renderer) rewindElement() {
	r.textPosition--
}
