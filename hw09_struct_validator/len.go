package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

type LenValidationError struct {
	Actual   int
	Expected int
}

func (e LenValidationError) Error() string {
	return fmt.Sprintf("Length must be %d, actual:%d", e.Expected, e.Actual)
}

func ValidateLen(value interface{}, params []string) (err error) {
	if len(params) != 1 || params[0] == "" {
		return NewValidatorError("length parameter is required")
	}
	lenExp, err := strconv.Atoi(params[0])
	if err != nil {
		return NewValidatorError("", err)
	}

	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch {
	case v.Kind() == reflect.String:
		lenActual := len([]rune(v.String()))
		if lenExp != lenActual {
			return LenValidationError{
				Actual:   lenActual,
				Expected: lenExp,
			}
		}
	case v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String:
		for i := 0; i < v.Len(); i++ {
			lenActual := len([]rune(v.Index(i).String()))
			if lenExp != lenActual {
				return LenValidationError{
					Actual:   lenActual,
					Expected: lenExp,
				}
			}
		}
	default:
		return NewValidatorError(fmt.Sprintf("incompatible type: %T", value))
	}

	return nil
}
