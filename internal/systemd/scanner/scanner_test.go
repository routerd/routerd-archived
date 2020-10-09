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
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"routerd.net/routerd/internal/systemd/token"
)

func TestScanner(t *testing.T) {
	// 	buf := bytes.NewBufferString(`[Match]
	// Name=enp1s0

	// [Network]
	// Description= test1 \
	// # address 1
	// test2 \
	// test3
	// # address 1
	// Address=10.1.10.9/24
	// # address 2
	// Address=10.1.10.11/24
	// Gateway=10.1.10.1 # test123
	// DNS=10.1.10.1

	// # route1000
	// # also important
	// [Route]
	// Gateway=192.168.0.11
	// Destination=10.0.0.0/8

	// # route2000
	// # this is very important!
	// [Route]
	// Gateway=192.168.0.12
	// Destination=20.0.0.0/8`)

	buf := bytes.NewBufferString(`[Match`)

	var s Scanner
	err := s.Init(buf)
	require.NoError(t, err)

	for {
		pos, tok, lit := s.Scan()
		fmt.Printf("%s\t\t%s\t%q\n", pos, tok, lit)
		if tok == token.EOF {
			break
		}
	}
	t.Fail()
}
