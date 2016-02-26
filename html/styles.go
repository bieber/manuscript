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

const inlineStyle = `
body {
	font-size: 20px;
}

div.container {
	width: 800px;
	margin-left: auto;
	margin-right: auto;
}

div.author_info {
	font-family: monospace;
	font-size: 12px;
}

h1 {
	font-size: 48px;
	text-align: center;
}

p.byline {
	text-align: center;
}

p.word_count {
	text-align: center;
}

div.short_story {
	position: relative;
}

div.short_story p.word_count {
	display: block;
	position: absolute;
	top: 0px;
	right: 0px;
}

div.table_of_contents {
	background-color: #eeeeee;
	display: inline-block;
	padding: 4px 16px 4px 4px;
}

div.table_of_contents ol {
	list-style: square;
}

div.table_of_contents li ol {
	list-style: disc;
}

div.scene {
	border-bottom: 2px solid #eeeeee;
}

h2 {
	text-align: center;
	font-size: 36px;
}

h3 {
	font-size: 28px;
}

p {
	text-indent: 60px;
}

div.front_matter p {
	text-indent: 0px;
}
`
