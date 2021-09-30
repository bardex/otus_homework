package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRegexp(t *testing.T) {
	text := "мир"
	var regexpError RegexpValidationError
	var valError ValidatorError
	var nilP *string
	var nilInter interface{}

	tests := []struct {
		text   interface{}
		regexp []string
		error  error
	}{
		// valid strings
		{text: nil, regexp: []string{"^мир$"}, error: nil},
		{text: nilInter, regexp: []string{"^мир$"}, error: nil},
		{text: nilP, regexp: []string{"^мир$"}, error: nil},
		{text: "мир", regexp: []string{"^мир$"}, error: nil},
		{text: "мир", regexp: []string{"(?i)^МИР$"}, error: nil},
		{text: "мировой", regexp: []string{"мир"}, error: nil},
		{text: UserRole("смирный"), regexp: []string{"мир"}, error: nil},
		{text: &text, regexp: []string{"мир"}, error: nil},
		{text: []string{"мир", "смирно"}, regexp: []string{"мир"}, error: nil},
		{text: []UserRole{UserRole("мир"), UserRole("мирный")}, regexp: []string{"^мир"}, error: nil},
		// invalid strings
		{text: "мир", regexp: []string{"^hello$"}, error: regexpError},
		{text: UserRole("смирный"), regexp: []string{"^мир$"}, error: regexpError},
		{text: []string{"hello", "смирно"}, regexp: []string{"hello"}, error: regexpError},
		{text: []UserRole{UserRole("hello"), UserRole("мир")}, regexp: []string{"hello"}, error: regexpError},
		// validator errors
		{text: 100, regexp: []string{"^hello$"}, error: valError},
		{text: []int{10}, regexp: []string{"^hello$"}, error: valError},
		{text: "мир", regexp: []string{"[a"}, error: valError},
		{text: "мир", regexp: []string{}, error: valError},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%s match %s", test.text, test.regexp)
		t.Run(name, func(t *testing.T) {
			err := ValidateRegexp(test.text, test.regexp)
			if test.error == nil {
				require.Nil(t, err)
			} else {
				require.ErrorAs(t, err, &test.error)
			}
		})
	}
}
