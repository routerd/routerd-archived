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
	"testing"

	"github.com/stretchr/testify/assert"
)

const example1 = `[Network]
Description= test1 \
# in the middle
test2 \
test3
# address 1
Address=10.1.10.9/24
# address 2
Address=10.1.10.11/24
`

var example1tokens = []tokenEntry{
	{tok: SECTION, lit: "[Network]"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Description"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: " test1 \\"},
	{tok: NEWLINE},
	{tok: COMMENT, lit: "# in the middle"},
	{tok: NEWLINE},
	{tok: STRING, lit: "test2 \\"},
	{tok: NEWLINE},
	{tok: STRING, lit: "test3"},
	{tok: NEWLINE},
	{tok: COMMENT, lit: "# address 1"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Address"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "10.1.10.9/24"},
	{tok: NEWLINE},
	{tok: COMMENT, lit: "# address 2"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Address"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "10.1.10.11/24"},
	{tok: NEWLINE},
	{tok: EOF},
}

const example2 = `# route1000
# also important
[Route]
Gateway=192.168.0.11
Destination=10.0.0.0/8

# route2000
# this is very important!
[Route]
Gateway=192.168.0.12
Destination=20.0.0.0/8`

var example2tokens = []tokenEntry{
	{tok: COMMENT, lit: "# route1000"},
	{tok: NEWLINE},
	{tok: COMMENT, lit: "# also important"},
	{tok: NEWLINE},
	{tok: SECTION, lit: "[Route]"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Gateway"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "192.168.0.11"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Destination"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "10.0.0.0/8"},
	{tok: NEWLINE},
	{tok: NEWLINE},
	{tok: COMMENT, lit: "# route2000"},
	{tok: NEWLINE},
	{tok: COMMENT, lit: "# this is very important!"},
	{tok: NEWLINE},
	{tok: SECTION, lit: "[Route]"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Gateway"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "192.168.0.12"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Destination"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "20.0.0.0/8"},
	{tok: EOF},
}

const example3 = `[Service]
Environment=ETCD_CA_FILE=/path/to/CA.pem
Environment=ETCD_CERT_FILE=/path/to/server.crt
Environment=ETCD_KEY_FILE=/path/to/server.key`

var example3tokens = []tokenEntry{
	{tok: SECTION, lit: "[Service]"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Environment"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "ETCD_CA_FILE"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "/path/to/CA.pem"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Environment"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "ETCD_CERT_FILE"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "/path/to/server.crt"},
	{tok: NEWLINE},
	{tok: STRING, lit: "Environment"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "ETCD_KEY_FILE"},
	{tok: ASSIGN}, // =
	{tok: STRING, lit: "/path/to/server.key"},
	{tok: EOF},
}

type tokenEntry struct {
	tok Token
	lit string
}

func TestScanner(t *testing.T) {
	tests := []struct {
		Name   string
		Input  string
		Tokens []tokenEntry
	}{
		{
			Name:   "Example 1",
			Input:  example1,
			Tokens: example1tokens,
		},
		{
			Name:   "Example 2",
			Input:  example2,
			Tokens: example2tokens,
		},
		{
			Name:   "Example 3",
			Input:  example3,
			Tokens: example3tokens,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Init
			var s Scanner
			s.Init([]byte(test.Input), nil)

			// Scan
			tokens := []tokenEntry{}
			for {
				pos, tok, lit := s.Scan()
				t.Logf("%s\t\t%s\t%q\n", pos, tok, lit)
				tokens = append(tokens, tokenEntry{
					tok: tok,
					lit: lit,
				})
				if tok == EOF {
					break
				}
			}
			assert.Equal(t, test.Tokens, tokens)
		})
	}
}
