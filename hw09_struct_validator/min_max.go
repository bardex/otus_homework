package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

type MinValidationError struct {
	Actual string
	Min    string
}

func (e MinValidationError) Error() string {
	return fmt.Sprintf("Min value is %s, actual:%s", e.Min, e.Actual)
}

type MaxValidationError struct {
	Actual string
	Max    string
}

func (e MaxValidationError) Error() string {
	return fmt.Sprintf("Max value is %s, actual:%s", e.Max, e.Actual)
}

type comparison struct {
	Value       string
	CompareWith string
	IsLess      bool
	IsMore      bool
}

func ValidateMin(value interface{}, params []string) (err error) {
	if len(params) != 1 {
		return NewValidatorError("min param is required")
	}
	if value == nil {
		return nil
	}

	cmp, err := compareNumbers(value, params[0])
	if err != nil {
		return err
	}
	if cmp.IsLess {
		return MinValidationError{
			Actual: cmp.Value,
			Min:    cmp.CompareWith,
		}
	}

	return nil
}

func ValidateMax(value interface{}, params []string) (err error) {
	if len(params) != 1 {
		return NewValidatorError("max param is required")
	}
	if value == nil {
		return nil
	}

	cmp, err := compareNumbers(value, params[0])
	if err != nil {
		return err
	}
	if cmp.IsMore {
		return MaxValidationError{
			Actual: cmp.Value,
			Max:    cmp.CompareWith,
		}
	}

	return nil
}

func compareNumbers(value interface{}, compareWith string) (*comparison, error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return &comparison{}, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		cmp, err := strconv.ParseInt(compareWith, 0, 64)
		if err != nil {
			return nil, NewValidatorError("", err)
		}
		value64 := v.Int()
		return &comparison{
			Value:       strconv.FormatInt(value64, 10),
			CompareWith: compareWith,
			IsLess:      value64 < cmp,
			IsMore:      value64 > cmp,
		}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		cmp, err := strconv.ParseUint(compareWith, 0, 64)
		if err != nil {
			return nil, NewValidatorError("", err)
		}
		value64 := v.Uint()
		return &comparison{
			Value:       strconv.FormatUint(value64, 10),
			CompareWith: compareWith,
			IsLess:      value64 < cmp,
			IsMore:      value64 > cmp,
		}, nil

	case reflect.Float32, reflect.Float64:
		cmp, err := strconv.ParseFloat(compareWith, 64)
		if err != nil {
			return nil, NewValidatorError("", err)
		}
		value64 := v.Float()
		return &comparison{
			Value:       strconv.FormatFloat(value64, 'g', 6, 64),
			CompareWith: compareWith,
			IsLess:      value64 < cmp,
			IsMore:      value64 > cmp,
		}, nil
	}

	return nil, NewValidatorError(fmt.Sprintf("incompatible type: %T", value))
}
