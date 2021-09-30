package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateLen(t *testing.T) {
	text := "мир!"
	empty := ""
	ctype := UserRole("мир!")
	var nilPtr *string
	var lenError LenValidationError
	var validatorError ValidatorError

	tests := []struct {
		value interface{}
		rule  string
		error error
	}{
		// nil
		{nil, "4", nil},
		{nilPtr, "4", nil},
		// valid strings
		{"мир!", "4", nil},
		{"", "0", nil},
		// invalid strings
		{"мир!", "5", lenError},
		{"", "1", lenError},
		// string pointers
		{&text, "4", nil},
		{&empty, "0", nil},
		{&empty, "1", lenError},
		// slice of strings
		{[]string{"пир", "мир"}, "3", nil},
		{[]string{"привет", "мир"}, "3", lenError},
		// custom type from string
		{UserRole("мир!"), "4", nil},
		{UserRole("мир!"), "2", lenError},
		// slice of custom types from string
		{[]UserRole{UserRole("пир"), UserRole("мир")}, "3", nil},
		{[]UserRole{UserRole("привет"), UserRole("мир")}, "3", lenError},
		// custom type pointer
		{&ctype, "4", nil},
		// syntax error
		{"мир!", "---", validatorError},
		// incompatible types
		{100, "3", validatorError},
		{[]byte{10, 20}, "2", validatorError},
		{[]rune{10, 20}, "2", validatorError},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v %v", test.value, test.rule)
		t.Run(name, func(t *testing.T) {
			err := ValidateLen(test.value, []string{test.rule})
			if test.error == nil {
				require.Nil(t, err)
			} else {
				require.ErrorAs(t, err, &test.error)
			}
		})
	}
}
