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

package renderers

import (
	"fmt"
	"github.com/bieber/manuscript/parser"
	"io"
	"regexp"
	"strings"
)

// RendererConstructor is a function that creates a new Renderer given
// a document and a set of options as string key/value pairs.
type RendererConstructor func(
	parser.Document,
	map[string]string,
) (Renderer, error)

// Renderer defines an object capable of rendering a document.
// to the given output.
type Renderer interface {
	Render(io.Writer) error
}

// Resolve attempts to find a match for the given document and
// renderer option string given the available set of renderer
// constructors.  If successful, it returns the newly instantiated
// renderer.
func Resolve(
	allRenderers map[string]RendererConstructor,
	document parser.Document,
	renderOption string,
) (Renderer, error) {
	matcher := regexp.MustCompile(
		`^(\w+)(?:\((\s*\w+\s*=\s*\w+\s*(?:,\s*\w+\s*=\s*\w+\s*)*)\))?$`,
	)
	matches := matcher.FindStringSubmatch(renderOption)
	if len(matches) != 3 {
		return nil, fmt.Errorf("Invalid renderer string %s", renderOption)
	}

	rendererName := matches[1]
	rendererArgs := map[string]string{}
	if matches[2] != "" {
		argSets := strings.Split(matches[2], ",")
		for _, argSet := range argSets {
			parts := strings.Split(argSet, "=")
			k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			rendererArgs[k] = v
		}
	}

	if constructor, ok := allRenderers[rendererName]; ok {
		return constructor(document, rendererArgs)
	}
	return nil, fmt.Errorf("%s is not a valid renderer", rendererName)
}
