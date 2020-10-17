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

type File struct {
	Sections []Section
}

func (f *File) SectionsByName(name string) (out []Section) {
	for _, section := range f.Sections {
		if section.Name != name {
			continue
		}

		out = append(out, section)
	}
	return
}

type Section struct {
	Name    string
	Comment string
	Keys    []Key
}

func (s *Section) KeysByName(name string) (out []Key) {
	for _, key := range s.Keys {
		if key.Name != name {
			continue
		}

		out = append(out, key)
	}
	return
}

type Key struct {
	Name    string
	Value   string
	Comment string
}
