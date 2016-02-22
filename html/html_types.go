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
)

type document struct {
	XMLName xml.Name `xml:"html"`
	Head    header
	Body    body
}

type header struct {
	XMLName    xml.Name `xml:"head"`
	Title      string   `xml:"title"`
	StyleSheet *link
}

type body struct {
	XMLName xml.Name `xml:"body"`
	Content div
}

type link struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
	HREF    string   `xml:"href,attr"`
}

type div struct {
	XMLName  xml.Name `xml:"div"`
	Class    string   `xml:"class,attr"`
	Children []interface{}
}

type h1 struct {
	XMLName xml.Name `xml:"h1"`
	Title   string   `xml:",chardata"`
}

type p struct {
	XMLName xml.Name `xml:"p"`
	Class   string   `xml:"class,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type span struct {
	XMLName xml.Name `xml:"span"`
	Class   string   `xml:"class,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type br struct {
	XMLName xml.Name `xml:"br"`
}
