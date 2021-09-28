package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

type InValidationError struct {
	Actual string
	In     []string
}

func (e InValidationError) Error() string {
	return fmt.Sprintf("Value must be in [%v], actual:%s", e.In, e.Actual)
}

func ValidateIn(value interface{}, params []string) (err error) {
	if len(params) == 0 {
		return NewValidatorError("array is required")
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

	var valStr string

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valStr = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valStr = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		valStr = strconv.FormatFloat(v.Float(), 'g', 6, 64)
	case reflect.String:
		valStr = v.String()
	default:
		return NewValidatorError(fmt.Sprintf("incompatible type: %T", value))
	}

	for _, param := range params {
		if param == valStr {
			return nil
		}
	}

	return InValidationError{
		Actual: valStr,
		In:     params,
	}
}
