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
	"unicode"
)

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
	Text []DocumentElement
}

type StoryType int

const (
	ShortStory StoryType = iota
	Novel
)

// This isn't great, performance-wise, but all my elements are just
// strings with different semantic information attached to the type
// name, so simple typeswitches should work for now.
type DocumentElement interface{}

type ParagraphBreak bool
type SceneBreak bool
type PlainText string
type ItalicText string
type BoldText string
type BoldItalicText string

func Parse(rawFIN io.Reader) (d Document, err error) {
	fin := bufio.NewReader(rawFIN)

	name, args := "", []string{}
	for name != "begin" {
		name, args, err = parseDirective(fin)
		if err != nil {
			return
		}

		switch name {
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

func parseDirective(fin *bufio.Reader) (name string, args []string, err error) {
	eatWhitespace(fin)

	r, _, err := fin.ReadRune()
	if r != '@' {
		err = errors.New("Expected directive")
		return
	}

	name, err = readWord(fin)
	if err != nil {
		return
	}

	for name != "begin" {
		err = eatWhitespace(fin)
		if err != nil {
			return
		}

		r, _, err = fin.ReadRune()
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
