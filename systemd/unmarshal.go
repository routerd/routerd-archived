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
	"reflect"
)

// Unmarshal parses the systemd unit data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	file, err := Decode(data)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	return unmarshalSections(file, rv)
}

func unmarshalSections(file *File, rv reflect.Value) error {
	// must be a pointer
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{rv.Type()}
	}

	knownSections := map[string]struct{}{}
	tv := rv.Elem().Type()
	for i := 0; i < rv.Elem().NumField(); i++ {
		f := rv.Elem().Field(i)
		sectionName := nameFromStructField(tv.Field(i))
		knownSections[sectionName] = struct{}{}
		sections := file.SectionsByName(sectionName)
		if len(sections) == 0 {
			// no section with this name
			// TODO load into "catch all property"
			continue
		}

		for _, section := range sections {
			var newObject reflect.Value
			switch f.Type().Kind() {
			case reflect.Struct:
				newObject = reflect.New(f.Type()).Elem()

			case reflect.Ptr, reflect.Slice:
				newObject = reflect.New(f.Type().Elem())
			}

			newObjectPtr := newObject
			if newObject.Kind() != reflect.Ptr {
				newObjectPtr = newObject.Addr()
			}

			if err := unmarshalKeys(&section, newObjectPtr); err != nil {
				return err
			}
			commentField := newObjectPtr.Elem().FieldByName("Comment")
			if commentField.IsValid() {
				commentField.Set(reflect.ValueOf(section.Comment))
			}

			if f.Type().Kind() != reflect.Slice {
				f.Set(newObject)
				continue
			}
			f.Set(reflect.Append(f, newObjectPtr.Elem()))
		}
	}

	// Add Sections that don't fit into any other place
	// if there is a AddSection function implemented.
	addSection := rv.MethodByName("AddSection")
	if !addSection.IsValid() {
		return nil
	}
	for _, section := range file.Sections {
		if _, ok := knownSections[section.Name]; ok {
			continue
		}
		addSection.Call([]reflect.Value{reflect.ValueOf(section)})
	}
	return nil
}

func unmarshalKeys(section *Section, rv reflect.Value) error {
	// must be a pointer
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{rv.Type()}
	}

	knownKeys := map[string]struct{}{}
	tv := rv.Elem().Type()
	for i := 0; i < rv.Elem().NumField(); i++ {
		f := rv.Elem().Field(i)

		keyName := nameFromStructField(tv.Field(i))
		knownKeys[keyName] = struct{}{}
		keys := section.KeysByName(keyName)
		if len(keys) == 0 {
			// no key with this name
			continue
		}

		var comment string
		switch f.Type().Kind() {
		case reflect.String:
			key := keys[len(keys)-1]
			f.SetString(key.Value)
			comment = key.Comment

		case reflect.Slice:
			if f.Type().Elem().Kind() != reflect.String {
				// wrong key type
				continue
			}

			var values []string
			for _, key := range keys {
				values = append(values, key.Value)
				if comment != "" {
					comment += "\n"
				}
				comment += key.Comment
			}
			f.Set(reflect.ValueOf(values))
		}

		// comment handling
		addComment := rv.MethodByName("AddKeyComment")
		if !addComment.IsValid() {
			continue
		}
		addComment.Call([]reflect.Value{
			reflect.ValueOf(keyName),
			reflect.ValueOf(comment),
		})
	}

	// Add Keys that don't fit into any other place
	// if there is a AddKey function implemented.
	addKey := rv.MethodByName("AddKey")
	if !addKey.IsValid() {
		return nil
	}
	for _, key := range section.Keys {
		if _, ok := knownKeys[key.Name]; ok {
			continue
		}
		addKey.Call([]reflect.Value{reflect.ValueOf(key)})
	}
	return nil
}

const fieldTagName = "systemd"

// returns the key or section name for a field
func nameFromStructField(structField reflect.StructField) string {
	if tag := structField.Tag.Get(fieldTagName); tag != "" {
		return tag
	}
	return structField.Name
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}
