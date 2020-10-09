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

// Package token defines constants representing the lexical tokens of a systemd ini file.
package token

import (
	"fmt"
	"strconv"
)

// Token is the type for lexical tokens of a systemd ini file.
type Token int

// List of Tokens.
const (
	// Special
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Values
	IDENT

	// Operators and delimiters
	NEWLINE // \n
	ASSIGN  // =
	LBRACK  // [
	RBRACK  // ]
)

var tokens = [...]string{
	// Special
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	// Values
	IDENT: "IDENT",

	// Operators and delimiters
	NEWLINE: "NEWLINE",
	ASSIGN:  "=",
	LBRACK:  "[",
	RBRACK:  "]",
}

func IsDelimiter(ch rune) bool {
	return ch == '\n' || ch == '=' || ch == '[' || ch == ']'
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

type Position struct {
	Line, Column int
}

func (pos *Position) IsValid() bool { return pos.Line > 0 }

func (pos Position) String() string {
	if !pos.IsValid() {
		return "-"
	}

	s := fmt.Sprintf("%d", pos.Line)
	if pos.Column != 0 {
		s += fmt.Sprintf(":%d", pos.Column)
	}
	return s
}
