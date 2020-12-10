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

package pdf

import (
	"fmt"
	"github.com/bieber/manuscript/parser"
	"github.com/bieber/manuscript/renderers"
	"github.com/bieber/manuscript/util"
	"github.com/dustin/go-humanize"
	"github.com/jung-kurt/gofpdf"
	"io"
	"strings"
)

const fontFamily = "Courier"

const ptsPerInch = 72
const fontSize = 12
const singleSpace = fontSize * 1.15
const doubleSpace = fontSize * 2

// Renderer provides a Render method to render the given document to a
// PDF file.
type Renderer struct {
	pageSize        string
	pageOrientation string
	document        parser.Document
	pdf             *gofpdf.Fpdf
}

// New creates a new Renderer given a document and options.
func New(
	document parser.Document,
	options map[string]string,
) (renderers.Renderer, error) {
	pageSize := "Letter"
	pageOrientation := "P"

	for k, v := range options {
		switch k {
		case "pageSize":
			pageSize = v
		case "pageOrientation":
			pageOrientation = v
		default:
			return nil, fmt.Errorf("Invalid PDF option %s", k)
		}
	}

	return &Renderer{
		pageSize:        pageSize,
		pageOrientation: pageOrientation,
		document:        document,
	}, nil
}

// Render writes the requested document out to the specified io.Writer
// as a PDF file formatted in manuscript format.
func (r *Renderer) Render(fout io.Writer) error {
	r.pdf = gofpdf.New(r.pageOrientation, "pt", r.pageSize, "")
	r.pdf.SetMargins(ptsPerInch, ptsPerInch, ptsPerInch)
	r.pdf.SetAutoPageBreak(true, ptsPerInch)
	r.pdf.SetHeaderFunc(r.writeHeader)
	r.pdf.AddPage()

	r.writeTitle()

	firstPart := true
	for _, p := range r.document.Parts {
		r.renderPart(p, firstPart)
		firstPart = false
	}

	return r.pdf.Output(fout)
}

func (r *Renderer) writeTitle() {
	pdf, document := r.pdf, r.document
	pdf.SetFont(fontFamily, "", fontSize)
	pdf.SetXY(ptsPerInch, ptsPerInch)

	authorBlockLines := []string{}
	if document.Author.Name != "" {
		authorBlockLines = append(authorBlockLines, document.Author.Name)
	}
	if len(document.Author.Address) != 0 {
		authorBlockLines = append(authorBlockLines, document.Author.Address...)
	}
	if document.Author.PhoneNumber != "" {
		authorBlockLines = append(authorBlockLines, document.Author.PhoneNumber)
	}
	if document.Author.EmailAddress != "" {
		authorBlockLines = append(
			authorBlockLines,
			document.Author.EmailAddress,
		)
	}
	if len(document.Author.ProfessionalOrgs) != 0 {
		authorBlockLines = append(authorBlockLines, "")
		authorBlockLines = append(
			authorBlockLines,
			document.Author.ProfessionalOrgs...,
		)
	}
	pdf.Write(singleSpace, strings.Join(authorBlockLines, "\n"))

	w, h := pdf.GetPageSize()
	byline := "by " + document.Author.Byline
	if document.Type == parser.Novel {
		byline = "a novel " + byline
	}

	pdf.SetXY(ptsPerInch, h/2)
	pdf.WriteAligned(
		w-2*ptsPerInch,
		singleSpace,
		document.Title,
		"C",
	)

	pdf.SetXY(ptsPerInch, h/2+doubleSpace)
	pdf.WriteAligned(
		w-2*ptsPerInch,
		singleSpace,
		byline,
		"C",
	)

	words := "about " + humanize.Comma(document.WordCount()) + " words"
	if document.Type == parser.ShortStory {
		pdf.SetXY(ptsPerInch, ptsPerInch)
		pdf.WriteAligned(
			// This calculation continues to baffle me, and I suspect that
			// there's something screwy going on in the gofpdf library.
			// For some reason using what seems like the appropriate width
			// (w - 2 * ptsPerInch) makes the header render too far away
			// from the right margin, but leaving out the -10 factor for
			// whatever reason causes it to line break even for very short
			// text.
			w-ptsPerInch-10,
			singleSpace,
			words,
			"R",
		)
		pdf.SetXY(2*ptsPerInch, h/2+4*doubleSpace)
	} else if document.Type == parser.Novel {
		pdf.SetXY(ptsPerInch, h-ptsPerInch-singleSpace)
		pdf.WriteAligned(
			w-2*ptsPerInch,
			singleSpace,
			words,
			"C",
		)
		pdf.SetX(2 * ptsPerInch)
	}
}

func (r *Renderer) renderPart(part parser.Part, firstInDocument bool) {
	pdf := r.pdf
	w, h := pdf.GetPageSize()
	if !part.Anonymous {
		text := util.PartLabel(part.Number, part.Title)
		pdf.AddPage()
		pdf.SetFont(fontFamily, "", fontSize)
		pdf.SetXY(ptsPerInch, h/2-2*doubleSpace)
		pdf.Bookmark(text, 0, -1)
		pdf.WriteAligned(
			w-2*ptsPerInch,
			singleSpace,
			text,
			"C",
		)
		pdf.SetXY(2*ptsPerInch, h/2)
	}

	firstChapter := !firstInDocument
	bookmarkLevel := 0
	if !part.Anonymous {
		bookmarkLevel++
	}
	for _, c := range part.Chapters {
		r.renderChapter(c, firstChapter, bookmarkLevel)
		firstChapter = false
	}
}

func (r *Renderer) renderChapter(
	chapter parser.Chapter,
	firstInPart bool,
	bookmarkLevel int,
) {
	pdf := r.pdf
	w, h := pdf.GetPageSize()

	if !chapter.Anonymous {
		if !firstInPart {
			pdf.AddPage()
		}
		pdf.SetFont(fontFamily, "", fontSize)
		pdf.SetXY(ptsPerInch, h/2)

		bookmarkText := ""
		labelText := ""
		if chapter.Prologue {
			bookmarkText = util.PrologueLabel(chapter.Title)
			labelText = "Prologue"
		} else {
			bookmarkText = util.ChapterLabel(chapter.Number, chapter.Title)
			labelText = fmt.Sprintf("Chapter %d", chapter.Number)
		}

		pdf.Bookmark(bookmarkText, bookmarkLevel, -1)
		pdf.WriteAligned(
			w-2*ptsPerInch,
			singleSpace,
			labelText,
			"C",
		)

		newY := h/2 + 2*doubleSpace
		if chapter.Title != "" {
			pdf.SetXY(ptsPerInch, h/2+doubleSpace)
			pdf.WriteAligned(
				w-2*ptsPerInch,
				singleSpace,
				chapter.Title,
				"C",
			)
			newY += doubleSpace
		}
		pdf.SetXY(2*ptsPerInch, newY)
	}

	for _, s := range chapter.Scenes {
		r.renderScene(s)
	}
}

func (r *Renderer) renderScene(scene parser.Scene) {
	pdf := r.pdf
	w, _ := pdf.GetPageSize()

	for _, p := range scene.Paragraphs {
		r.renderParagraph(p)
	}

	if scene.EndsWithSceneBreak {
		// This is another addition I don't fully understand.  Without
		// this line, Using WriteAligned at the very beginning of a
		// page seems to cause some bizarre linebreak behavior in the
		// header, but if I write a single space before the hash mark,
		// which doesn't seem to visibly affect the rendering, the
		// problem goes away.
		pdf.Write(singleSpace, " ")
		pdf.WriteAligned(w-2*ptsPerInch, doubleSpace, "#", "C")
		pdf.Write(doubleSpace, "\n")
		pdf.SetX(2 * ptsPerInch)
	}
}

func (r *Renderer) renderParagraph(paragraph parser.Paragraph) {
	pdf := r.pdf

	for _, element := range paragraph.Text {
		switch e := element.(type) {
		case parser.PlainText:
			pdf.SetFont(fontFamily, "", fontSize)
			pdf.Write(doubleSpace, string(e))

		case parser.ItalicText:
			pdf.SetFont(fontFamily, "U", fontSize)
			pdf.Write(doubleSpace, string(e))

		case parser.BoldText:
			pdf.SetFont(fontFamily, "B", fontSize)
			pdf.Write(doubleSpace, string(e))

		case parser.BoldItalicText:
			pdf.SetFont(fontFamily, "BU", fontSize)
			pdf.Write(doubleSpace, string(e))

		}
	}

	pdf.Write(doubleSpace, "\n")
	pdf.SetX(2 * ptsPerInch)
}

func (r *Renderer) writeHeader() {
	pdf, document := r.pdf, r.document
	if pdf.PageNo() == 1 {
		return
	}

	pageNumber := pdf.PageNo()
	if document.Type == parser.Novel {
		pageNumber--
	}

	w, _ := pdf.GetPageSize()
	pdf.SetXY(ptsPerInch, ptsPerInch)
	pdf.WriteAligned(
		// This calculation continues to baffle me, and I suspect that
		// there's something screwy going on in the gofpdf library.
		// For some reason using what seems like the appropriate width
		// (w - 2 * ptsPerInch) makes the header render too far away
		// from the right margin, but leaving out the -10 factor for
		// whatever reason causes it to line break even for very short
		// text.
		w-ptsPerInch-10,
		singleSpace,
		fmt.Sprintf(
			"%s / %s / %d",
			document.Author.ShortName,
			document.ShortTitle,
			pageNumber,
		),
		"R",
	)
	pdf.SetXY(ptsPerInch, ptsPerInch+doubleSpace)
}
