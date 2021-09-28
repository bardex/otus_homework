package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateIn(t *testing.T) {
	text := "мир"
	in := 10
	var ptrText *string
	var ptrInt *int
	var emptyStr *string
	var emptyInt *int
	var inError InValidationError
	var valError ValidatorError

	ptrText = &text
	ptrInt = &in

	tests := []struct {
		value interface{}
		in    []string
		error error
	}{
		// valid values
		{value: nil, in: []string{"", "2"}, error: nil},
		{value: "мир", in: []string{"", "мир", "2"}, error: nil},
		{value: "", in: []string{"мир", ""}, error: nil},
		{value: 132, in: []string{"", "0", "132"}, error: nil},
		{value: uint8(100), in: []string{"", "0", "100"}, error: nil},
		{value: 0, in: []string{"", "0", "100"}, error: nil},
		{value: "0", in: []string{"0"}, error: nil},
		{value: "", in: []string{""}, error: nil},
		{value: ptrText, in: []string{"", "мир"}, error: nil},
		{value: ptrInt, in: []string{"0", "", "10"}, error: nil},
		{value: emptyStr, in: []string{"0", "10"}, error: nil},
		{value: emptyInt, in: []string{"0", "10"}, error: nil},
		{value: 100.5, in: []string{"0", "100.0", "100.5"}, error: nil},
		{value: 100.0, in: []string{"0", "100"}, error: nil},
		// invalid values
		{value: "мир", in: []string{"", "2"}, error: inError},
		{value: "", in: []string{"0"}, error: inError},
		{value: 132, in: []string{"10"}, error: inError},
		{value: uint8(100), in: []string{"", "0"}, error: inError},
		{value: 0, in: []string{"", "100"}, error: inError},
		{value: "0", in: []string{""}, error: inError},
		{value: ptrText, in: []string{"", "hello"}, error: inError},
		{value: ptrInt, in: []string{"0", "", "100"}, error: inError},
		{value: 100.5, in: []string{"0", "100", "100.51"}, error: inError},
		{value: 100.0, in: []string{"0", "100.1"}, error: inError},
		// validator errors
		{value: "мир", in: []string{}, error: valError},
		{value: true, in: []string{"0"}, error: valError},
		{value: []int{10}, in: []string{"10"}, error: valError},
		{value: []string{"10"}, in: []string{"10"}, error: valError},
		{value: 10 + 5i, in: []string{"10+5i"}, error: valError},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%v in %v", test.value, test.in)
		t.Run(name, func(t *testing.T) {
			err := ValidateIn(test.value, test.in)
			if test.error == nil {
				require.Nil(t, err)
			} else {
				require.ErrorAs(t, err, &test.error)
			}
		})
	}
}
