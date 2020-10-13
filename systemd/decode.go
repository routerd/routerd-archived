/*
Copyright 2020 The routerd Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package systemd

import (
	"fmt"
	"strings"

	"routerd.net/routerd/systemd/internal"
)

// Decode takes a systemd configuration file and returns a data container to access and manipulate it.
func Decode(data []byte) (*File, error) {
	var d decodeState
	d.init(data)
	return d.decode()
}

// decodeState stores the current state of a decode operation.
type decodeState struct {
	scanner internal.Scanner
	// comments belonging to the next section
	// or the next/current key
	comment string
	section *Section // active section
	key     *Key     // active key
	file    *File
}

func (d *decodeState) init(src []byte) *decodeState {
	d.scanner.Init(src, nil)
	d.comment = ""
	d.section = nil
	d.key = nil
	d.file = &File{}
	return d
}

func (d *decodeState) decode() (*File, error) {
decode:
	for {
		pos, tok, lit := d.scanner.Scan()
		switch tok {
		case internal.COMMENT:
			d.addComment(pos, tok, lit)

		case internal.ASSIGN:
			if d.key != nil {
				// ignore ASSIGN tokens
				// when scanning a value
				d.key.Value += "="
			}

		case internal.EOF:
			// force close
			d.closeKey()
			break decode

		case internal.NEWLINE:
			if d.key != nil &&
				!strings.HasSuffix(d.key.Value, "\\") {
				// stop scanning value, but continue scanning on \
				// for multi line strings
				d.closeKey()
			}

		case internal.STRING:
			if err := d.addString(pos, tok, lit); err != nil {
				return nil, err
			}

		case internal.SECTION:
			if err := d.addSection(pos, tok, lit); err != nil {
				return nil, err
			}
		}
	}
	return d.file, nil
}

func (d *decodeState) addSection(pos internal.Position, tok internal.Token, lit string) error {
	// validate section name
	if !strings.HasPrefix(lit, "[") {
		return fmt.Errorf("%s: section needs to start with [, is: %q", pos, lit)
	}
	if !strings.HasSuffix(lit, "]") {
		return fmt.Errorf("%s: section needs to end with ], is: %q", pos, lit)
	}

	d.file.Sections = append(d.file.Sections, Section{
		Name:    lit[1 : len(lit)-1], // strip [ ]
		Comment: d.comment,
	})
	d.section = &d.file.Sections[len(d.file.Sections)-1]
	d.comment = ""
	return nil
}

func (d *decodeState) addString(pos internal.Position, tok internal.Token, lit string) error {
	if d.section == nil {
		// We want to be in a section before encountering any STRING
		return fmt.Errorf("%s: key started outside of section %q", pos, lit)
	}

	// KEY
	if d.key == nil {
		if pos, tok, lit := d.scanner.Scan(); tok != internal.ASSIGN {
			return fmt.Errorf("%s: key not followed by = (ASSIGN), token found: %s %q", pos, tok, lit)
		}

		d.section.Keys = append(d.section.Keys, Key{
			Name: strings.TrimSpace(lit),
		})
		d.key = &d.section.Keys[len(d.section.Keys)-1]
		return nil
	}

	// Value
	d.key.Value += strings.TrimSpace(lit)
	return nil
}

func (d *decodeState) addComment(pos internal.Position, tok internal.Token, lit string) {
	if d.comment != "" {
		d.comment += "\n"
	}
	d.comment += strings.TrimSpace(lit[1:]) // strip # or ;
}

func (d *decodeState) closeKey() {
	if d.key == nil {
		return
	}
	d.key.Comment = d.comment
	d.key.Value = strings.ReplaceAll(d.key.Value, "\\", " ")
	d.key = nil
	d.comment = ""
}
