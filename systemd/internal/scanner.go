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

package internal

import (
	"unicode/utf8"
)

// An ErrorHandler may be provided to Scanner.Init.
type ErrorHandler func(pos Position, msg string)

// Scanner implements a scanner for systemd unit files.
// It takes a []byte as source which can then be tokenized
// through repeated calls to the Scan method.
type Scanner struct {
	pos Position
	src []byte
	err ErrorHandler

	// scanning state
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset (position after current character)

	ErrorCount int // number of errors encountered
}

func (s *Scanner) Init(src []byte, err ErrorHandler) {
	s.src = src
	s.pos.Line = 1
	s.pos.Column = 1

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.ErrorCount = 0

	s.next()
}

func (s *Scanner) error(msg string) {
	if s.err != nil {
		s.err(s.pos, msg)
	}
	s.ErrorCount++
}

func (s *Scanner) next() {
	if s.rdOffset >= len(s.src) {
		// we are at the end of our buffer
		// -> EOF
		s.offset = len(s.src)
		s.ch = -1
		return
	}

	s.offset = s.rdOffset
	if s.ch == '\n' {
		s.pos.Line++
		s.pos.Column = 1
	} else {
		s.pos.Column++
	}

	r, w := rune(s.src[s.rdOffset]), 1
	switch {
	case r == 0:
		s.error("illegal character NUL")
	case r >= utf8.RuneSelf:
		r, w = utf8.DecodeRune(s.src[s.rdOffset:])
		if r == utf8.RuneError && w == 1 {
			s.error("illegal UTF-8 encoding")
		}
	}
	s.ch = r
	s.rdOffset += w
}

func (s *Scanner) scanString() string {
	offs := s.offset - 1
	for !IsDelimiter(s.ch) && s.ch != -1 {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanSection() string {
	offs := s.offset - 1
	for s.ch != ']' && s.ch != '\n' && s.ch != -1 {
		s.next()
	}
	if s.ch == ']' {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) Scan() (pos Position, tok Token, lit string) {
	pos = s.pos

	ch := s.ch
	s.next()
	switch ch {

	case -1:
		tok = EOF

	case '\n':
		s.pos.Line++
		s.pos.Column = 0
		tok = NEWLINE

	case '[':
		tok = SECTION
		lit = s.scanSection()

	case '#', ';':
		tok = COMMENT
		lit = s.scanString()

	case '=':
		tok = ASSIGN

	default:
		tok = STRING
		lit = s.scanString()
	}
	return
}
