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

// KeyComments implements a comment container to allow
// attachment of comments to keys.
type KeyComments struct {
	comments map[string]string
}

func (c *KeyComments) GetKeyComment(key string) string {
	return c.comments[key]
}

func (c *KeyComments) AddKeyComment(key, comment string) {
	if c.comments == nil {
		c.comments = map[string]string{}
	}
	c.comments[key] = comment
}

func (c *KeyComments) RemoveKeyComment(key string) {
	if c.comments != nil {
		delete(c.comments, key)
	}
}

// SectionList is storing arbitrary sections.
// When embedded into a struct that is given to Unmarshal,
// sections that cannot be assigned to a field are
// appended to this SectionList
type SectionList []Section

func (l *SectionList) AddSection(s Section) {
	*l = append(*l, s)
}

// KeyList is storing arbitrary key.
// When embedded into a struct that is given to Unmarshal,
// keys that cannot be assigned to a field are
// appended to this KeyList
type KeyList []Key

func (l *KeyList) AddKey(k Key) {
	*l = append(*l, k)
}
