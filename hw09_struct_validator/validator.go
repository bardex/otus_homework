package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidatorError struct {
	Msg string
	Err error
}

func (e ValidatorError) Unwrap() error {
	return e.Err
}

func (e ValidatorError) Error() string {
	return e.Msg
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	for _, err := range v {
		b.WriteString(fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error()))
	}
	return b.String()
}

type validator func(value interface{}, params []string) error

var validators = map[string]validator{
	"len":    ValidateLen,
	"min":    ValidateMin,
	"max":    ValidateMax,
	"in":     ValidateIn,
	"regexp": ValidateRegexp,
}

func Validate(v interface{}) error {
	r := reflect.ValueOf(v)

	if r.Kind() != reflect.Struct {
		return NewValidatorError(fmt.Sprintf("%T is not struct", v))
	}

	invalids := make(ValidationErrors, 0)

	for i := 0; i < r.NumField(); i++ {
		if !r.Field(i).CanInterface() {
			continue
		}

		fType := r.Type().Field(i)
		fTag := fType.Tag.Get("validate")
		fName := fType.Name
		fValue := r.Field(i).Interface()

		rules, err := parseValidateTag(fTag)
		if err != nil {
			return NewValidatorError(fmt.Sprintf("`%s` has validator error:", fName), err)
		}

		for _, rule := range rules {
			validator, ok := validators[rule.Rule]
			if !ok {
				return NewValidatorError(fmt.Sprintf("`%s` unknown validator: `%s`", fName, rule.Rule))
			}
			if err := validator(fValue, rule.Params); err != nil {
				if errors.As(err, &ValidatorError{}) {
					return NewValidatorError(fmt.Sprintf("`%s` has validator error: `%s`:", fName, rule.Rule), err)
				}
				invalids = append(invalids, ValidationError{Field: fName, Err: err})
			}
		}
	}

	if len(invalids) > 0 {
		return invalids
	}

	return nil
}

func NewValidatorError(msg string, err ...error) ValidatorError {
	e := ValidatorError{Msg: msg}
	if len(err) == 1 {
		e.Err = err[0]
	}
	return e
}
