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

package scanner

import (
	"io"
	"io/ioutil"
	"unicode/utf8"

	"routerd.net/routerd/internal/systemd/token"
)

type ErrorHandler func(pos token.Position, msg string)

type Scanner struct {
	pos token.Position
	src []byte
	err ErrorHandler

	// scanning state
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset (position after current character)

	ErrorCount int // number of errors encountered
}

func (s *Scanner) Init(reader io.Reader) error {
	var err error
	s.src, err = ioutil.ReadAll(reader)
	s.pos.Line = 1
	s.pos.Column = 1

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.ErrorCount = 0

	s.next()
	return err
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

func (s *Scanner) scanIdentifier() string {
	offs := s.offset - 1
	for !token.IsDelimiter(s.ch) && s.ch != -1 {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) Scan() (pos token.Position, tok token.Token, lit string) {
	pos = s.pos

	ch := s.ch
	s.next()
	switch ch {

	case -1:
		tok = token.EOF

	case '\n':
		s.pos.Line++
		s.pos.Column = 0
		tok = token.NEWLINE

	case '[':
		tok = token.LBRACK

	case ']':
		tok = token.RBRACK

	case '#', ';':
		tok = token.COMMENT
		lit = s.scanIdentifier()

	case '=':
		tok = token.ASSIGN

	default:
		tok = token.IDENT
		lit = s.scanIdentifier()
	}
	return
}
