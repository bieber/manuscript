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
	"github.com/StefanSchroeder/Golang-Roman"
	"github.com/bieber/manuscript/parser"
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

// Render writes the requested document out to the specified io.Writer
// as a PDF file formatted in manuscript format.
func Render(fout io.Writer, document parser.Document) error {
	pdf := gofpdf.New("P", "pt", "Letter", "")
	pdf.SetMargins(ptsPerInch, ptsPerInch, ptsPerInch)
	pdf.SetAutoPageBreak(true, ptsPerInch)
	pdf.SetHeaderFunc(func() { writeHeader(pdf, document) })
	pdf.AddPage()

	writeTitle(pdf, document)

	partNumber := 0
	chapterNumber := 0
	var lastElement parser.DocumentElement
	for _, e := range document.Text {
		switch e.(type) {
		case parser.PrologueBreak:
			pdf.AddPage()

		case parser.PartBreak:
			partNumber++
			chapterNumber = 0
			pdf.AddPage()

		case parser.ChapterBreak:
			chapterNumber++
			if _, ok := lastElement.(parser.PartBreak); !ok {
				pdf.AddPage()
			}
		}
		writeElement(pdf, e, partNumber, chapterNumber)
		lastElement = e
	}

	return pdf.Output(fout)
}

func writeTitle(pdf *gofpdf.Fpdf, document parser.Document) {
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

func writeElement(
	pdf *gofpdf.Fpdf,
	element parser.DocumentElement,
	partNumber int,
	chapterNumber int,
) {
	w, h := pdf.GetPageSize()

	switch e := element.(type) {
	case parser.ParagraphBreak:
		pdf.Write(doubleSpace, "\n")
		pdf.SetX(2 * ptsPerInch)

	case parser.PrologueBreak:
		bookmarkText := "Prologue"
		if e != "" {
			bookmarkText = bookmarkText + ": " + string(e)
		}

		pdf.SetFont(fontFamily, "", fontSize)
		pdf.SetXY(ptsPerInch, h/2)
		pdf.Bookmark(bookmarkText, 0, -1)

		pdf.WriteAligned(
			w-2*ptsPerInch,
			singleSpace,
			"Prologue",
			"C",
		)

		newY := h/2 + 2*doubleSpace
		if e != "" {
			pdf.SetXY(ptsPerInch, h/2+doubleSpace)
			pdf.WriteAligned(
				w-2*ptsPerInch,
				singleSpace,
				string(e),
				"C",
			)
			newY += doubleSpace
		}
		pdf.SetXY(2*ptsPerInch, newY)

	case parser.PartBreak:
		text := "Part " + roman.Roman(partNumber)
		if e != "" {
			text += ": " + string(e)
		}

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

	case parser.ChapterBreak:
		bookmarkText := fmt.Sprintf("Chapter %d", chapterNumber)
		if e != "" {
			bookmarkText = bookmarkText + ": " + string(e)
		}
		bookmarkLevel := 0
		if partNumber != 0 {
			bookmarkLevel = 1
		}

		pdf.SetFont(fontFamily, "", fontSize)
		pdf.SetXY(ptsPerInch, h/2)
		pdf.Bookmark(bookmarkText, bookmarkLevel, -1)
		pdf.WriteAligned(
			w-2*ptsPerInch,
			singleSpace,
			fmt.Sprintf("Chapter %d", chapterNumber),
			"C",
		)

		newY := h/2 + 2*doubleSpace
		if e != "" {
			pdf.SetXY(ptsPerInch, h/2+doubleSpace)
			pdf.WriteAligned(
				w-2*ptsPerInch,
				singleSpace,
				string(e),
				"C",
			)
			newY += doubleSpace
		}
		pdf.SetXY(2*ptsPerInch, newY)

	case parser.SceneBreak:

		pdf.WriteAligned(w-2*ptsPerInch, doubleSpace, "#", "C")
		pdf.Write(doubleSpace, "\n")
		pdf.SetX(2 * ptsPerInch)

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

func writeHeader(pdf *gofpdf.Fpdf, document parser.Document) {
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
