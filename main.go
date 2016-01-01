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

package main

import (
	"fmt"
	"github.com/bieber/conflag"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

type Config struct {
	Help     bool
	Renderer string
	Output   string
}

func main() {
	config := &Config{
		Renderer: "pdf",
	}

	configParser, err := conflag.New(config)
	if err != nil {
		log.Fatal(err)
	}

	configParser.ProgramName("manuscript")
	configParser.ProgramDescription("" +
		"Usage: manuscript (-o | --output) outfile [options] infile\n\n" +
		"Format stories in manuscript format.  For input format details, see " +
		"README file.",
	)

	configParser.Field("Help").
		ShortFlag('h').
		LongFlag("help").
		Description("Print usage text and exit.")
	configParser.Field("Renderer").
		ShortFlag('r').
		LongFlag("renderer").
		Description("Select a renderer for your story.")
	configParser.Field("Output").
		ShortFlag('o').
		LongFlag("output").
		Required().
		Description("File path to write output to.")
	configParser.AllowExtraArgs("input")

	extraArgs, err := configParser.Read()
	if err != nil || len(extraArgs) != 1 || config.Help {
		exitCode := 0

		if err != nil {
			log.Println(err)
			exitCode = 1
		}

		if width, _, err := terminal.GetSize(0); err == nil {
			fmt.Println(configParser.Usage(uint(width)))
		}

		os.Exit(exitCode)
	}
}