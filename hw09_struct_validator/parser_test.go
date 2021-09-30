package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseValidateTag(t *testing.T) {
	var validatorError ValidatorError

	tests := []struct {
		tag   string
		rules validationRules
		error error
	}{
		{tag: "", rules: validationRules{}, error: nil},
		{
			tag:   "required",
			error: nil,
			rules: validationRules{
				validationRule{Rule: "required", Params: []string{}},
			},
		},
		{
			tag:   "max:50",
			error: nil,
			rules: validationRules{
				validationRule{Rule: "max", Params: []string{"50"}},
			},
		},
		{
			tag:   "min:0|max:50",
			error: nil,
			rules: validationRules{
				validationRule{Rule: "min", Params: []string{"0"}},
				validationRule{Rule: "max", Params: []string{"50"}},
			},
		},
		{
			tag:   "min :0|required |max:20| in: 1 , 2 ,3",
			error: nil,
			rules: validationRules{
				validationRule{Rule: "min", Params: []string{"0"}},
				validationRule{Rule: "required", Params: []string{}},
				validationRule{Rule: "max", Params: []string{"20"}},
				validationRule{Rule: "in", Params: []string{"1", "2", "3"}},
			},
		},
		{tag: "  ", rules: nil, error: validatorError},
		{tag: "|", rules: nil, error: validatorError},
		{tag: "required||min:0", rules: nil, error: validatorError},
		{tag: ":", rules: nil, error: validatorError},
		{tag: "max:", rules: nil, error: validatorError},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			rules, err := parseValidateTag(test.tag)
			if test.error == nil {
				require.Nil(t, err)
			} else {
				require.ErrorAs(t, err, &test.error)
			}
			require.Equal(t, test.rules, rules)
		})
	}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := parseValidateTag("min :0|required |max:20| in: 1 , 2 ,3")
		if err != nil {
			b.Fatal(err)
		}
	}
}
