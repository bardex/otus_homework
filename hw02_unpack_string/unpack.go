package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	stateNormal = iota
	stateEscape
	escaper = '\\'
)

func Unpack(str string) (string, error) {
	var buffer rune
	var result strings.Builder

	flushBuffer := func(repeat int) {
		if buffer > 0 {
			result.WriteString(strings.Repeat(string(buffer), repeat))
		}
		buffer = 0
	}

	state := stateNormal
	for _, char := range str {
		switch {
		case state == stateNormal && char == escaper:
			flushBuffer(1)
			state = stateEscape

		case state == stateNormal && unicode.IsDigit(char):
			if buffer == 0 {
				return "", ErrInvalidString
			}
			repeat, err := strconv.Atoi(string(char))
			if err != nil {
				return "", err
			}
			flushBuffer(repeat)

		case state == stateNormal:
			flushBuffer(1)
			buffer = char

		case state == stateEscape && (unicode.IsDigit(char) || char == escaper):
			buffer = char
			state = stateNormal

		default:
			return "", ErrInvalidString
		}
	}

	// проход по строке должен заканчиваться в нормальном состоянии т.е. не должно не примененных слэшей или цифр
	if state != stateNormal {
		return "", ErrInvalidString
	}

	// последний
	flushBuffer(1)
	return result.String(), nil
}
