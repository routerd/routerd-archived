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

	"routerd.net/routerd/systemd/internal"
)

// func Parse(src []byte) *File {
// 	var scanner internal.Scanner
// 	scanner.Init(src, nil)

// 	var (
// 		section *Section
// 		key     *Key
// 	)
// 	for {
// 		_, tok, lit := s.Scan()
// 		if tok == EOF {
// 			break
// 		}

// 		if key != nil {
// 			// we are scanning a key value
// 			continue
// 		}

// 		if section != nil {
// 			// we are within a section
// 			continue
// 		}

// 		// we need to search for a section

// 	}
// }

type parser struct {
	scanner internal.Scanner
	// comments []string // comment lines
	// section  *Section // active section
	// key      *Key     // active key
}

func (p *parser) Init(src []byte) {
	p.scanner.Init(src, nil)
}

func (p *parser) Parse() (*File, error) {
	var (
		comments []string // comment lines
		section  *Section // active section
		key      *Key     // active key
	)

	for {
		_, tok, lit := p.scanner.Scan()
		if tok == internal.EOF {
			break
		}

		if tok == internal.COMMENT {
			comments = append(comments, lit)
			continue
		}

		if key != nil {
			// we are scanning a key value
			continue
		}

		if section != nil {
			// we are within a section
			continue
		}

		// we need to search for a section
		if tok != internal.SECTION {
			// whatever is not starting a new section...
			return nil, fmt.Errorf("unexpected token %s %q", tok, lit)
		}

	}
	return nil, nil
}

// func (p *parser.)
