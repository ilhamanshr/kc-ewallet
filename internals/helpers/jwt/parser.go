package jwt

import (
	"reflect"
)

// Parsing claim permission to []string output
func ParsePermission(input interface{}, output *[]string) {
	var res []string
	s := reflect.ValueOf(input)

	if s.IsValid() {
		for i := 0; i < s.Len(); i++ {
			res = append(res, s.Index(i).Interface().(string))
		}
	}

	*output = res
}
