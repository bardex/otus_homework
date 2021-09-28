package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
)

type RegexpValidationError struct {
	Actual string
	Regexp string
}

func (e RegexpValidationError) Error() string {
	return fmt.Sprintf("the string must match to `%s`, actual:%s", e.Regexp, e.Actual)
}

// use lru cache in production.
var (
	regexpCache      = map[string]*regexp.Regexp{}
	regexpCacheMutex sync.Mutex
)

func ValidateRegexp(value interface{}, params []string) error {
	if len(params) != 1 {
		return NewValidatorError("regexp param is required")
	}
	if value == nil {
		return nil
	}

	re, err := makeRegexp(params[0])
	if err != nil {
		return err
	}

	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.String {
		if re.MatchString(v.String()) {
			return nil
		}
		return RegexpValidationError{
			Actual: v.String(),
			Regexp: params[0],
		}
	}

	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String {
		for i := 0; i < v.Len(); i++ {
			if !re.MatchString(v.Index(i).String()) {
				return RegexpValidationError{
					Actual: v.Index(i).String(),
					Regexp: params[0],
				}
			}
		}
		return nil
	}

	return NewValidatorError(fmt.Sprintf("incompatible type: %T", value))
}

func makeRegexp(rexp string) (*regexp.Regexp, error) {
	regexpCacheMutex.Lock()
	defer regexpCacheMutex.Unlock()

	re, exists := regexpCache[rexp]
	if !exists {
		var err error
		if re, err = regexp.Compile(rexp); err != nil {
			return nil, NewValidatorError("", err)
		}
		regexpCache[rexp] = re
	}
	return re, nil
}
