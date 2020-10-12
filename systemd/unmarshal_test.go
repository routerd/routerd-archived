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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFile struct {
	SectionList
	Match   *matchSection
	Network networkSection
	Routes  []routeSection `systemd:"Route"`
}

type matchSection struct {
	KeyComments
	KeyList
	Comment string
	Name    string
}

type networkSection struct {
	Addresses []string `systemd:"Address"`
}

type routeSection struct {
	KeyList
	Gateway     string
	Destination string
}

func TestUnmarshal(t *testing.T) {
	f := &testFile{}
	err := Unmarshal([]byte(`# this is a config file!
[Match]
# some comment
# more comment!
Name=eth*

[Network]
Address=10.10.10.2/24
Address=10.10.10.3/24

# a section comment!
[Route]
Gateway=10.10.10.1/24
# comment for dest key
Destination=10.10.20.1/24

[Route]
Gateway=10.10.10.1/24
UndefinedKey=something

[Whatever]
`), f)

	expected := &testFile{
		Match: &matchSection{
			KeyComments: KeyComments{
				comments: map[string]string{
					"Name": "some comment\nmore comment!",
				},
			},
			Comment: "this is a config file!",
			Name:    "eth*",
		},
		Network: networkSection{
			Addresses: []string{
				"10.10.10.2/24",
				"10.10.10.3/24",
			},
		},
		Routes: []routeSection{
			{
				Gateway:     "10.10.10.1/24",
				Destination: "10.10.20.1/24",
			},
			{
				Gateway: "10.10.10.1/24",
				KeyList: KeyList{
					{Name: "UndefinedKey", Value: "something"},
				},
			},
		},
		SectionList: SectionList{
			{Name: "Whatever"},
		},
	}

	require.NoError(t, err)
	assert.Equal(t, expected, f)
}
