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

type Section struct {
	Name    string
	Comment string
	Keys    []Key
}

type Key struct {
	Name    string
	Value   string
	Comment string
}

// func (f *File) Section(section string) Section {
// 	for _, s := range f.Sections {
// 		if s.Name == section {
// 			return s
// 		}
// 	}
// 	return Section{}
// }

// func (s *Section) Key(key string) Key {
// 	for _, s := range f.Sections {
// 		if s.Name == section {
// 			return s
// 		}
// 	}
// 	return Section{}
// }
