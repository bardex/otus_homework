package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\3`, expected: `3`},
		{input: `\32`, expected: `33`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "abc14", `qw\ne`, `\`, `a3\`, `a\`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			s, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q, result %s", err, s)
		})
	}
}

func TestMultiByteUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "♬2Ё4丕2", expected: "♬♬ЁЁЁЁ丕丕"},
		{input: "♬0丕", expected: "丕"},
		{input: `丕\2Ё`, expected: "丕2Ё"},
		{input: `♬0Ё0`, expected: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}
