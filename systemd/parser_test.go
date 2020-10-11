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

const example1 = `# network comment
[Network]
# start desc
Description= test1 \
	# in the middle
	test2 \
	test3
# address 1
Address=10.1.10.9/24

# address 2
Address=10.1.10.11/24
`

var example1File = &File{
	Sections: []Section{
		{
			Comment: "network comment",
			Name:    "Network",
			Keys: []Key{
				{
					Comment: "start desc\nin the middle",
					Name:    "Description",
					Value:   "test1  test2  test3",
				},
				{
					Comment: "address 1",
					Name:    "Address",
					Value:   "10.1.10.9/24",
				},
				{
					Comment: "address 2",
					Name:    "Address",
					Value:   "10.1.10.11/24",
				},
			},
		},
	},
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

var example2File = &File{
	Sections: []Section{
		{
			Comment: "route1000\nalso important",
			Name:    "Route",
			Keys: []Key{
				{
					Name:  "Gateway",
					Value: "192.168.0.11",
				},
				{
					Name:  "Destination",
					Value: "10.0.0.0/8",
				},
			},
		},
		{
			Comment: "route2000\nthis is very important!",
			Name:    "Route",
			Keys: []Key{
				{
					Name:  "Gateway",
					Value: "192.168.0.12",
				},
				{
					Name:  "Destination",
					Value: "20.0.0.0/8",
				},
			},
		},
	},
}

const example3 = `[Service]
Environment=ETCD_CA_FILE=/path/to/CA.pem
Environment=ETCD_CERT_FILE=/path/to/server.crt
Environment=ETCD_KEY_FILE=/path/to/server.key`

var example3File = &File{
	Sections: []Section{
		{
			Name: "Service",
			Keys: []Key{
				{
					Name:  "Environment",
					Value: "ETCD_CA_FILE=/path/to/CA.pem",
				},
				{
					Name:  "Environment",
					Value: "ETCD_CERT_FILE=/path/to/server.crt",
				},
				{
					Name:  "Environment",
					Value: "ETCD_KEY_FILE=/path/to/server.key",
				},
			},
		},
	},
}

func TestParser(t *testing.T) {
	tests := []struct {
		Name  string
		Input string
		File  *File
	}{
		{
			Name:  "example1",
			Input: example1,
			File:  example1File,
		},
		{
			Name:  "example2",
			Input: example2,
			File:  example2File,
		},
		{
			Name:  "example3",
			Input: example3,
			File:  example3File,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var d decodeState
			f, err := d.init([]byte(test.Input)).unmarshal()

			require.NoError(t, err)
			assert.Equal(t, test.File, f)
		})
	}
}
