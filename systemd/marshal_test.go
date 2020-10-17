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

func TestMarshal(t *testing.T) {
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
				Enable:      BoolPtr(true),
			},
			{
				Gateway: "10.10.10.1/24",
				Source:  StringPtr("something"),
				KeyList: KeyList{
					{Name: "UndefinedKey", Value: "something"},
				},
			},
		},
		SectionList: SectionList{
			{Name: "Whatever"},
		},
	}

	b, err := Marshal(expected)
	require.NoError(t, err)

	assert.Equal(t, `# this is a config file!
[Match]
# some comment
# more comment!
Name=eth*

[Network]
Address=10.10.10.2/24
Address=10.10.10.3/24

[Route]
Gateway=10.10.10.1/24
Destination=10.10.20.1/24
Enable=yes
Disable=

[Route]
Gateway=10.10.10.1/24
Source=something
Disable=
UndefinedKey=something

[Whatever]
`, string(b))
}
