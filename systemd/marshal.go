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
	"bytes"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)

	// must be a pointer
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, &InvalidUnmarshalError{rv.Type()}
	}

	file := &File{}
	tv := rv.Elem().Type()
	for i := 0; i < rv.Elem().NumField(); i++ {
		structField := tv.Field(i)
		if structField.Type.Name() == "SectionList" {
			continue
		}

		field := rv.Elem().Field(i)
		sectionName := nameFromStructField(structField)
		switch field.Type().Kind() {
		case reflect.Ptr:
			section := Section{
				Name: sectionName,
			}
			marshalSection(&section, field)
			file.Sections = append(file.Sections, section)

		case reflect.Struct:
			section := Section{
				Name: sectionName,
			}
			marshalSection(&section, field.Addr())
			file.Sections = append(file.Sections, section)

		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				section := Section{
					Name: sectionName,
				}
				marshalSection(&section, field.Index(i).Addr())
				file.Sections = append(file.Sections, section)
			}
		}

	}

	// Check, if the SectionList for arbitrary sections is embedded.
	sectionList := rv.Elem().FieldByName("SectionList")
	if sectionList.IsValid() {
		for i := 0; i < sectionList.Len(); i++ {
			file.Sections = append(file.Sections, sectionList.Index(i).Interface().(Section))
		}
	}

	// Encode to bytes
	var out bytes.Buffer
	if err := Encode(&out, file); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func marshalSection(section *Section, rv reflect.Value) {
	tv := rv.Elem().Type()
	for i := 0; i < rv.Elem().NumField(); i++ {
		structField := tv.Field(i)
		field := rv.Elem().Field(i)

		if structField.Name == "Comment" &&
			structField.Type.Kind() == reflect.String {
			section.Comment = field.String()
			continue
		}
		if structField.Name == "KeyComments" ||
			structField.Name == "KeyList" {
			continue
		}

		keyName := nameFromStructField(structField)
		switch structField.Type.Kind() {
		case reflect.String:
			key := Key{
				Name:    keyName,
				Value:   field.String(),
				Comment: keyComment(rv, keyName),
			}
			if key.Value == "" {
				continue
			}

			section.Keys = append(section.Keys, key)

		case reflect.Slice:
			for i, val := range field.Interface().([]string) {
				key := Key{
					Name:  keyName,
					Value: val,
				}
				if key.Value == "" {
					continue
				}

				if i == 0 {
					// Add the comment to the first Key
					key.Comment = keyComment(rv, keyName)
				}
				section.Keys = append(section.Keys, key)
			}
		}
	}

	// Check if KeyList for arbitrary keys is embedded.
	keyList := rv.Elem().FieldByName("KeyList")
	if !keyList.IsValid() {
		return
	}
	for i := 0; i < keyList.Len(); i++ {
		section.Keys = append(section.Keys, keyList.Index(i).Interface().(Key))
	}
}

// keyComment gets the registered comment for the given key
func keyComment(rv reflect.Value, keyName string) string {
	getComment := rv.MethodByName("GetKeyComment")
	if getComment.IsValid() {
		comment := getComment.Call([]reflect.Value{
			reflect.ValueOf(keyName),
		})
		if len(comment) == 1 {
			return comment[0].String()
		}
	}
	return ""
}
