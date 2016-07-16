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
	"errors"
	"fmt"
	"github.com/StefanSchroeder/Golang-Roman"
	"github.com/bieber/manuscript/parser"
	"github.com/bieber/manuscript/renderers"
	"github.com/dustin/go-humanize"
	"io"
	"strings"
)

// Renderer provides a Render method to render the given document to
// an HTML file.
type Renderer struct {
	styleSheet string
	authorInfo bool
	includeTOC bool
	document   parser.Document
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

	return &renderer, nil
}

// Render writes the requested document out to the specified io.Writer
// as an HTML file.
func (r *Renderer) Render(fout io.Writer) error {
	encoder := xml.NewEncoder(selfClosingRemover{fout})

	bodyContents := []interface{}{}
	bodyContents = append(bodyContents, r.renderFrontMatter())

	if r.includeTOC {
		toc := r.renderTOC()
		if len(toc.Children) != 0 {
			bodyContents = append(bodyContents, toc)
		}
	}

	for _, p := range r.document.Parts {
		bodyContents = append(bodyContents, r.renderPart(p))
	}

	storyTypeClass := ""
	if r.document.Type == parser.Novel {
		storyTypeClass = " novel"
	} else if r.document.Type == parser.ShortStory {
		storyTypeClass = " short_story"
	}

	encoder.Indent("", "\t")
	return encoder.Encode(
		document{
			Head: r.renderHead(),
			Body: body{
				Content: div{
					Class:    "container" + storyTypeClass,
					Children: bodyContents,
				},
			},
		},
	)
}

func (r *Renderer) renderHead() header {
	var styleSheet *link
	var inlineStyleSheet *style
	if r.styleSheet == "" {
		rawStyle := inlineStyle

		styleLines := strings.Split(rawStyle, "\n")
		for i := range styleLines {
			if i != len(styleLines)-1 {
				styleLines[i] = "\t\t\t" + styleLines[i]
			}
		}

		inlineStyleSheet = &style{Text: strings.Join(styleLines, "\n") + "\t\t"}
	} else if r.styleSheet != "" {
		styleSheet = &link{
			Rel:  "stylesheet",
			Type: "text/css",
			HREF: r.styleSheet,
		}
	}

	return header{
		Title:      r.document.Title,
		StyleSheet: styleSheet,
		Style:      inlineStyleSheet,
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
	outerChildren := []interface{}{}

	for _, p := range r.document.Parts {
		children := []interface{}{}
		for _, c := range p.Chapters {
			if c.Anonymous {
				continue
			}

			text, href := "", ""
			if c.Prologue {
				text = "Prologue"
				if c.Title != "" {
					text += ": " + c.Title
				}
				href = fmt.Sprintf("#prologue_%d_%d", p.Number, c.Number)
			} else {
				text = fmt.Sprintf("Chapter %d", c.Number)
				if c.Title != "" {
					text += ": " + c.Title
				}
				href = fmt.Sprintf("#chapter_%d_%d", p.Number, c.Number)
			}

			children = append(
				children,
				li{
					Children: []interface{}{
						a{
							Text: text,
							HREF: href,
						},
					},
				},
			)
		}

		if len(children) == 0 {
			continue
		}

		if p.Anonymous {
			outerChildren = append(outerChildren, children...)
		} else {
			text := "Part " + roman.Roman(p.Number)
			if p.Title != "" {
				text += ": " + p.Title
			}

			outerChildren = append(
				outerChildren,
				li{
					Children: []interface{}{
						a{
							Text: text,
							HREF: fmt.Sprintf("#part_%d", p.Number),
						},
						ol{
							Children: children,
						},
					},
				},
			)
		}
	}

	if len(outerChildren) == 0 {
		return div{}
	}

	return div{
		Class: "table_of_contents",
		Children: []interface{}{
			ol{Class: "toc_outer", Children: outerChildren},
		},
	}
}

func (r *Renderer) renderPart(part parser.Part) div {
	class := "anonymous_part"
	children := []interface{}{}

	if !part.Anonymous {
		class = "part"
		text := "Part " + roman.Roman(part.Number)
		if part.Title != "" {
			text += ": " + part.Title
		}

		children = append(
			children,
			h2{
				Children: []interface{}{
					a{
						Name: fmt.Sprintf("part_%d", part.Number),
						Text: text,
					},
				},
			},
		)
	}

	for _, c := range part.Chapters {
		children = append(children, r.renderChapter(c, part.Number))
	}

	return div{
		Class:    class,
		Children: children,
	}

}

func (r *Renderer) renderChapter(chapter parser.Chapter, partNumber int) div {
	class := "anonymous_chapter"
	children := []interface{}{}

	if !chapter.Anonymous {
		if chapter.Prologue {
			class = "chapter prologue"

			text := "Prologue"
			if chapter.Title != "" {
				text += ": " + chapter.Title
			}

			children = append(
				children,
				h3{
					Children: []interface{}{
						a{
							Name: fmt.Sprintf(
								"prologue_%d_%d",
								partNumber,
								chapter.Number,
							),
							Text: text,
						},
					},
				},
			)
		} else {
			class = "chapter"

			text := fmt.Sprintf("Chapter %d", chapter.Number)
			if chapter.Title != "" {
				text += ": " + chapter.Title
			}

			children = append(
				children,
				h3{
					Children: []interface{}{
						a{
							Name: fmt.Sprintf(
								"chapter_%d_%d",
								partNumber,
								chapter.Number,
							),
							Text: text,
						},
					},
				},
			)
		}
	}

	for _, s := range chapter.Scenes {
		children = append(children, r.renderScene(s))
	}

	return div{
		Class:    class,
		Children: children,
	}
}

func (r *Renderer) renderScene(scene parser.Scene) div {
	children := []interface{}{}
	for _, p := range scene.Paragraphs {
		children = append(children, r.renderParagraph(p))
	}

	return div{
		Class:    "scene",
		Children: children,
	}
}

func (r *Renderer) renderParagraph(paragraph parser.Paragraph) p {
	children := []interface{}{}
	for _, e := range paragraph.Text {
		children = append(children, r.renderElement(e))
	}

	return p{Children: children}
}

func (r *Renderer) renderElement(element parser.DocumentElement) interface{} {
	switch e := element.(type) {
	case parser.PlainText:
		return span{Text: string(e)}
	case parser.ItalicText:
		return em{Text: string(e)}
	case parser.BoldText:
		return strong{Text: string(e)}
	case parser.BoldItalicText:
		return strong{Child: em{Text: string(e)}}
	default:
		panic(
			errors.New(
				"html: Unexpected document element passed to renderElement",
			),
		)
	}
}
