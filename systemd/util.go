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

// Returns a pointer to the given string.
func StringPtr(str string) *string {
	return &str
}

// Returns a pointer to the given bool.
func BoolPtr(b bool) *bool {
	return &b
}

func StrToBool(b string) *bool {
	switch b {
	case "1", "yes", "true", "on":
		return BoolPtr(true)

	case "0", "no", "false", "off":
		return BoolPtr(false)
	}
	return nil
}

func BoolToStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
