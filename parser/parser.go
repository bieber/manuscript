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
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

// Document defines a story, both its text and relevant metadata.
type Document struct {
	Type       StoryType
	Title      string
	ShortTitle string
	Author     struct {
		Name             string
		Byline           string
		ShortName        string
		Address          []string
		PhoneNumber      string
		EmailAddress     string
		ProfessionalOrgs []string
	}
	Parts []Part
	Text  []DocumentElement
}

// Part defines a part of the document, which may or may not have a
// title, and may also be anonymous (meaning that the document hasn't
// explicitly declared the beginning of a part and no title page
// should be emitted).
type Part struct {
	Title     string
	Anonymous bool
	Number    int

	Chapters []Chapter
}

// Chapter defines a chapter of the document, which may or may not
// have a title, and may also be anonymous (meaning that the document
// hasn't explicitly declared the beginning of a chapter and no title
// should be emitted).  A Chapter may also be a prologue, which is
// essentially the same but with a different type of header and
// doesn't contribute to chapter numbering.
type Chapter struct {
	Title     string
	Anonymous bool
	Prologue  bool
	Number    int

	Scenes []Scene
}

// Scene defines a single scene in the text, which may or may not end
// with a hard scene-break.
type Scene struct {
	EndsWithSceneBreak bool

	Paragraphs []Paragraph
}

// Paragraph defines a single paragraph of text, composed of
// potentially multiple sections of text with varying formatting.
type Paragraph struct {
	Text []DocumentElement
}

// StoryType defines the type of a document.
type StoryType int

const (
	// ShortStory is a story without parts or chapters.
	ShortStory StoryType = iota
	// Novel is a longer story which may have a prologue, parts and/or
	// chapters.
	Novel
)

// DocumentElement is just an empty interface.  I'm using type
// switches to differentiate between the different types of element.
type DocumentElement interface{}

// ParagraphBreak is just a linebreak between paragraphs.
type ParagraphBreak bool

// SceneBreak is a break between scenes.
type SceneBreak bool

// PrologueBreak is a break in the text for a prologue.  It may have a
// title or be empty.
type PrologueBreak string

// PartBreak is a break for a new part of a story.  It may have a
// title or be empty.
type PartBreak string

// ChapterBreak is a break for a new chapter in the story.  It may
// have a title or be empty.
type ChapterBreak string

// PlainText is simple unformatted text.
type PlainText string

// ItalicText will be rendered as italic.
type ItalicText string

// BoldText will be rendered as bold.
type BoldText string

// BoldItalicText will be rendered as both bold and italic.
type BoldItalicText string

// Parse reads a document from a text file and returns a parsed
// Document object if there aren't any errors.
func Parse(rawFIN io.Reader) (d Document, err error) {
	fin := bufio.NewReader(rawFIN)

	d, err = lexMetadata(fin)
	if err != nil {
		return
	}

	for {
		es := []DocumentElement{}
		es, err = lexParagraphOrDirective(fin)

		if err == io.EOF {
			if es != nil {
				d.Text = append(d.Text, es...)
			}
			err = nil

			d.Parts = parseText(d.Text)
			return
		}
		if err != nil {
			return
		}

		d.Text = append(d.Text, es...)
	}

	return
}

func lexMetadata(fin *bufio.Reader) (d Document, err error) {
	name, args := "", []string{}
	for name != "begin" {
		name, args, err = lexMetadataDirective(fin)
		if err != nil {
			return
		}

		switch name {
		case "notes":
			continue

		case "type":
			if len(args) != 1 {
				err = errors.New("Missing type")
				return
			}

			switch args[0] {
			case "shortStory":
				d.Type = ShortStory
			case "novel":
				d.Type = Novel
			default:
				err = errors.New("Invalid story type")
				return
			}

		case "title":
			if len(args) != 1 {
				err = errors.New("Missing title")
				return
			}
			d.Title = args[0]

		case "shortTitle":
			if len(args) != 1 {
				err = errors.New("Missing short title")
				return
			}
			d.ShortTitle = args[0]

		case "authorName":
			if len(args) != 1 {
				err = errors.New("Missing author name")
				return
			}
			d.Author.Name = args[0]

		case "authorShortName":
			if len(args) != 1 {
				err = errors.New("Missing author short name")
				return
			}
			d.Author.ShortName = args[0]

		case "authorByline":
			if len(args) != 1 {
				err = errors.New("Missing author byline")
				return
			}
			d.Author.Byline = args[0]

		case "authorAddress":
			if len(args) < 1 {
				err = errors.New("Missing author address")
				return
			}
			d.Author.Address = args

		case "authorPhoneNumber":
			if len(args) != 1 {
				err = errors.New("Missing author phone number")
				return
			}
			d.Author.PhoneNumber = args[0]

		case "authorEmail":
			if len(args) != 1 {
				err = errors.New("Missing author email")
				return
			}
			d.Author.EmailAddress = args[0]

		case "authorOrgs":
			if len(args) < 1 {
				err = errors.New("Missing author organizations")
				return
			}
			d.Author.ProfessionalOrgs = args

		case "begin":
			break

		default:
			err = errors.New("Unrecognized directive")
			return
		}
	}

	return
}

func lexParagraphOrDirective(
	fin *bufio.Reader,
) (es []DocumentElement, err error) {
	err = eatWhitespace(fin)
	if err != nil {
		return nil, err
	}

	r := '\000'
	r, _, err = fin.ReadRune()
	if err != nil {
		return
	}
	if r == '@' {
		fin.UnreadRune()

		var e DocumentElement
		e, err = lexDirective(fin)
		if err != nil {
			return
		}
		if e != nil {
			es = []DocumentElement{e}
		}
	} else {
		fin.UnreadRune()
		es, err = lexParagraph(fin)
	}

	return
}

// The key to metadata directives is that they will always be
// terminated by the beginning '@' of another directive (except for
// @begin), and their arguments may span multiple lines.
func lexMetadataDirective(
	fin *bufio.Reader,
) (name string, args []string, err error) {
	err = eatWhitespace(fin)
	if err != nil {
		return
	}

	r, _, err := fin.ReadRune()
	if r != '@' {
		err = errors.New("Expected directive")
		return
	}

	name, err = readWord(fin)
	if err != nil {
		return
	}

	for name != "begin" && name != "scene" {
		err = eatWhitespace(fin)
		if err != nil {
			return
		}

		r, _, err = fin.ReadRune()
		if err != nil {
			return
		}

		fin.UnreadRune()
		if r == '@' {
			break
		}

		arg := ""
		arg, err = readPlainText(fin)
		if err != nil {
			return
		}
		args = append(args, arg)
	}

	return
}

// A regular directive in the text may only have a single,
// newline-terminated argument.
func lexDirective(fin *bufio.Reader) (e DocumentElement, err error) {
	r := '\000'
	r, _, err = fin.ReadRune()
	if r != '@' {
		err = errors.New("Missing '@' in directive")
	}
	if err != nil {
		return
	}

	name := ""
	name, err = readWord(fin)
	if err != nil {
		return
	}

	argDirectives := map[string]bool{
		"chapter":  true,
		"part":     true,
		"prologue": true,
		"note":     true,
	}

	if name == "scene" {
		e = SceneBreak(true)
		return
	} else if _, ok := argDirectives[name]; !ok {
		err = errors.New("Invalid directive")
		return
	}

	rawArg := []rune{}
	for {
		r, _, err = fin.ReadRune()
		if err != nil {
			return
		}
		if r == '\n' {
			break
		}
		rawArg = append(rawArg, r)
	}
	arg := strings.TrimSpace(string(rawArg))

	if name == "chapter" {
		e = ChapterBreak(arg)
	} else if name == "part" {
		e = PartBreak(arg)
	} else if name == "prologue" {
		e = PrologueBreak(arg)
	}

	return
}

func lexParagraph(fin *bufio.Reader) (es []DocumentElement, err error) {
	buf := []rune{}
	bold := false
	italic := false

	for {
		r := '\000'
		r, _, err = fin.ReadRune()
		if err != nil {
			return
		}

		if r == '\n' {
			r, _, err = fin.ReadRune()
			if err != nil {
				if err == io.EOF {
					if len(buf) != 0 {
						es = append(es, formatText(buf, bold, italic))
					}
				}
				return
			}

			fin.UnreadRune()
			if r == '\n' || r == '@' {
				if len(buf) != 0 {
					es = append(es, formatText(buf, bold, italic))
				}
				break
			} else {
				buf = addWhitespace(buf)
			}
		} else if unicode.IsSpace(r) {
			buf = addWhitespace(buf)
		} else if r == '\\' {
			r, _, err = fin.ReadRune()
			if err != nil {
				return
			}
			buf = append(buf, r)
		} else if r == '*' {
			flipItalic := true
			flipBold := false

			r, _, err = fin.ReadRune()
			if err != nil {
				return
			}

			if r == '*' {
				flipBold = true
				flipItalic = false

				r, _, err = fin.ReadRune()
				if err != nil {
					return
				}

				if r == '*' {
					flipItalic = true
				} else {
					fin.UnreadRune()
				}
			} else {
				fin.UnreadRune()
			}

			es = append(es, formatText(buf, bold, italic))
			buf = []rune{}

			if flipBold {
				bold = !bold
			}
			if flipItalic {
				italic = !italic
			}
		} else {
			buf = append(buf, r)
		}
	}

	es = append(es, ParagraphBreak(true))
	return
}

func parseText(text []DocumentElement) (ps []Part) {
	var p Part
	for partNumber := 0; len(text) != 0; {
		p, text = parsePart(text)

		if !p.Anonymous {
			partNumber++
		}
		p.Number = partNumber

		ps = append(ps, p)
	}
	return
}

func parsePart(text []DocumentElement) (p Part, rest []DocumentElement) {
	if partBreak, ok := text[0].(PartBreak); ok {
		p.Anonymous = false
		p.Title = string(partBreak)
		text = text[1:]
	} else {
		p.Anonymous = true
	}

	var c Chapter
	for chapterNumber, prologueNumber := 0, 0; len(text) != 0; {
		c, text = parseChapter(text)

		if c.Prologue {
			if !c.Anonymous {
				prologueNumber++
			}
			c.Number = prologueNumber
		} else {
			if !c.Anonymous {
				chapterNumber++
			}
			c.Number = chapterNumber
		}

		p.Chapters = append(p.Chapters, c)
		if len(text) != 0 {
			if _, ok := text[0].(PartBreak); ok {
				break
			}
		}
	}

	rest = text
	return
}

func parseChapter(text []DocumentElement) (c Chapter, rest []DocumentElement) {
	if prologueBreak, ok := text[0].(PrologueBreak); ok {
		c.Anonymous = false
		c.Prologue = true
		c.Title = string(prologueBreak)
		text = text[1:]
	} else if chapterBreak, ok := text[0].(ChapterBreak); ok {
		c.Anonymous = false
		c.Title = string(chapterBreak)
		text = text[1:]
	} else {
		c.Anonymous = true
	}

	var s Scene
outer:
	for len(text) != 0 {
		s, text = parseScene(text)

		c.Scenes = append(c.Scenes, s)
		if len(text) != 0 {
			switch text[0].(type) {
			case PrologueBreak:
				break outer
			case ChapterBreak:
				break outer
			case PartBreak:
				break outer
			}
		}
	}

	rest = text
	return
}

func parseScene(text []DocumentElement) (s Scene, rest []DocumentElement) {
	var p Paragraph
outer:
	for len(text) != 0 {
		p, text = parseParagraph(text)

		s.Paragraphs = append(s.Paragraphs, p)
		if len(text) != 0 {
			switch text[0].(type) {
			case SceneBreak:
				text = text[1:]
				s.EndsWithSceneBreak = true
				break outer
			case PrologueBreak:
				break outer
			case ChapterBreak:
				break outer
			case PartBreak:
				break outer

			}
		}
	}

	rest = text
	return
}

func parseParagraph(
	text []DocumentElement,
) (p Paragraph, rest []DocumentElement) {
outer:
	for len(text) != 0 {
		switch text[0].(type) {
		case ParagraphBreak:
			text = text[1:]
			break outer
		case SceneBreak:
			break outer
		case PrologueBreak:
			break outer
		case ChapterBreak:
			break outer
		case PartBreak:
			break outer
		}

		p.Text = append(p.Text, text[0])
		text = text[1:]
	}

	rest = text
	return
}

func formatText(text []rune, bold, italic bool) DocumentElement {
	if italic && bold {
		return BoldItalicText(text)
	} else if bold {
		return BoldText(text)
	} else if italic {
		return ItalicText(text)
	}
	return PlainText(text)
}

func addWhitespace(text []rune) []rune {
	if len(text) == 0 || text[len(text)-1] != ' ' {
		text = append(text, ' ')
	}
	return text
}

func eatWhitespace(fin *bufio.Reader) error {
	for {
		r, _, err := fin.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			err = fin.UnreadRune()
			return nil
		}
	}
}

func readWord(fin *bufio.Reader) (text string, err error) {
	chars := []rune{}
	for {
		r := '\000'
		r, _, err = fin.ReadRune()
		if err != nil {
			return
		}

		if unicode.IsSpace(r) {
			fin.UnreadRune()
			break
		} else {
			chars = append(chars, r)
		}
	}

	text = string(chars)
	return
}

func readPlainText(fin *bufio.Reader) (text string, err error) {
	chars := []rune{}
	for {
		r := '\000'
		r, _, err = fin.ReadRune()
		if err != nil {
			return
		}

		if r == '\n' {
			break
		} else {
			chars = append(chars, r)
		}
	}

	text = string(chars)
	return
}
