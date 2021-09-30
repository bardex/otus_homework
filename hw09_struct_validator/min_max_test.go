package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestUserID int64

func TestValidateMin(t *testing.T) {
	number := 100
	ctype := TestUserID(100)
	var empty *int
	var minError MinValidationError
	var valError ValidatorError

	tests := []struct {
		value interface{}
		min   []string
		error error
	}{
		// nil
		{value: nil, min: []string{"100"}, error: nil},
		// valid numbers
		{value: 100, min: []string{"0"}, error: nil},
		{value: 100, min: []string{"100"}, error: nil},
		{value: -100, min: []string{"-100"}, error: nil},
		{value: int(101), min: []string{"100"}, error: nil},
		{value: int8(102), min: []string{"100"}, error: nil},
		{value: int16(1022), min: []string{"100"}, error: nil},
		{value: int32(103), min: []string{"100"}, error: nil},
		{value: int64(104), min: []string{"100"}, error: nil},
		{value: uint8(105), min: []string{"100"}, error: nil},
		{value: uint16(1055), min: []string{"100"}, error: nil},
		{value: uint32(106), min: []string{"100"}, error: nil},
		{value: uint64(107), min: []string{"100"}, error: nil},
		{value: TestUserID(108), min: []string{"100"}, error: nil},
		{value: &number, min: []string{"100"}, error: nil},
		{value: &ctype, min: []string{"100"}, error: nil},
		{value: empty, min: []string{"100"}, error: nil},
		{value: float32(2e3), min: []string{"1.5e3"}, error: nil},
		{value: float64(2e13), min: []string{"1.5e13"}, error: nil},
		{value: 0.75, min: []string{"0.75"}, error: nil},
		{value: 0.75, min: []string{"0.5"}, error: nil},
		{value: -0.75, min: []string{"-0.75"}, error: nil},
		{value: -0.75, min: []string{"-0.8"}, error: nil},
		{value: 10.0, min: []string{"10"}, error: nil},

		// invalid numbers
		{value: -200, min: []string{"0"}, error: minError},
		{value: 200, min: []string{"10000000000000"}, error: minError},
		{value: -200, min: []string{"-199"}, error: minError},
		{value: int(202), min: []string{"203"}, error: minError},
		{value: int8(20), min: []string{"21"}, error: minError},
		{value: int16(22), min: []string{"23"}, error: minError},
		{value: int32(203), min: []string{"300"}, error: minError},
		{value: int64(204), min: []string{"300"}, error: minError},
		{value: uint8(205), min: []string{"300"}, error: minError},
		{value: uint16(55), min: []string{"60"}, error: minError},
		{value: uint32(206), min: []string{"207"}, error: minError},
		{value: uint64(207), min: []string{"208"}, error: minError},
		{value: TestUserID(208), min: []string{"300"}, error: minError},
		{value: &number, min: []string{"200"}, error: minError},
		{value: &ctype, min: []string{"200"}, error: minError},
		{value: float32(1.5555e3), min: []string{"2e3"}, error: minError},

		// validator errors
		{value: 120, min: []string{}, error: valError},
		{value: 125, min: []string{"---"}, error: valError},
		{value: uint64(207), min: []string{"---"}, error: valError},
		{value: 130, min: []string{""}, error: valError},
		{value: 100.52, min: []string{"10w2"}, error: valError},
		{value: "мир!", min: []string{"100"}, error: valError},
		{value: UserRole("мир"), min: []string{"100"}, error: valError},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%v>=%v", test.value, test.min)
		t.Run(name, func(t *testing.T) {
			err := ValidateMin(test.value, test.min)
			if test.error == nil {
				require.Nil(t, err)
			} else {
				require.ErrorAs(t, err, &test.error)
			}
		})
	}
}

func TestValidateMax(t *testing.T) {
	number := 100
	ctype := TestUserID(100)
	var maxError MinValidationError
	var valError ValidatorError

	tests := []struct {
		value interface{}
		max   []string
		error error
	}{
		// nil
		{value: nil, max: []string{"100"}, error: nil},
		// valid numbers
		{value: -100, max: []string{"0"}, error: nil},
		{value: 100, max: []string{"110"}, error: nil},
		{value: -10000, max: []string{"-9000"}, error: nil},
		{value: int(101), max: []string{"101"}, error: nil},
		{value: int8(102), max: []string{"102"}, error: nil},
		{value: int16(1022), max: []string{"2000000000000000"}, error: nil},
		{value: int32(103), max: []string{"200"}, error: nil},
		{value: int64(104), max: []string{"105"}, error: nil},
		{value: uint8(105), max: []string{"105"}, error: nil},
		{value: uint16(155), max: []string{"156"}, error: nil},
		{value: uint32(106), max: []string{"107"}, error: nil},
		{value: uint64(107), max: []string{"107"}, error: nil},
		{value: TestUserID(108), max: []string{"200"}, error: nil},
		{value: &number, max: []string{"100"}, error: nil},
		{value: &ctype, max: []string{"100"}, error: nil},
		{value: float32(2e3), max: []string{"2.12e3"}, error: nil},
		{value: float64(2e13), max: []string{"2.012e13"}, error: nil},
		{value: 0.75, max: []string{"0.75"}, error: nil},
		{value: 0.75, max: []string{"0.77"}, error: nil},
		{value: -0.75, max: []string{"-0.55"}, error: nil},
		{value: -0.75, max: []string{"-0.750"}, error: nil},
		{value: 10.0, max: []string{"10"}, error: nil},

		// invalid numbers
		{value: 200, max: []string{"0"}, error: maxError},
		{value: 20000000000000, max: []string{"10"}, error: maxError},
		{value: -200, max: []string{"-300"}, error: maxError},
		{value: int(202), max: []string{"200"}, error: maxError},
		{value: int8(20), max: []string{"19"}, error: maxError},
		{value: int16(22), max: []string{"19"}, error: maxError},
		{value: int32(203), max: []string{"200"}, error: maxError},
		{value: int64(204), max: []string{"200"}, error: maxError},
		{value: uint8(205), max: []string{"200"}, error: maxError},
		{value: uint16(55), max: []string{"50"}, error: maxError},
		{value: uint32(206), max: []string{"200"}, error: maxError},
		{value: uint64(207), max: []string{"200"}, error: maxError},
		{value: TestUserID(208), max: []string{"200"}, error: maxError},
		{value: &number, max: []string{"10"}, error: maxError},
		{value: &ctype, max: []string{"10"}, error: maxError},
		{value: float32(1.5555e3), max: []string{"1e3"}, error: maxError},

		// validator errors
		{value: 120, max: []string{}, error: valError},
		{value: 125, max: []string{"---"}, error: valError},
		{value: uint64(207), max: []string{"---"}, error: valError},
		{value: 130, max: []string{""}, error: valError},
		{value: 100.52, max: []string{"10w2"}, error: valError},
		{value: "мир!", max: []string{"100"}, error: valError},
		{value: UserRole("мир"), max: []string{"100"}, error: valError},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%v<=%v", test.value, test.max)
		t.Run(name, func(t *testing.T) {
			err := ValidateMax(test.value, test.max)
			if test.error == nil {
				require.Nil(t, err)
			} else {
				require.ErrorAs(t, err, &test.error)
			}
		})
	}
}
