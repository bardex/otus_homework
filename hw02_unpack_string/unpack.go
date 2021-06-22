package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	esc := '\\'
	slice := []rune(str)
	last := len(slice) - 1

	for i := 0; i <= last; i++ {
		char := slice[i]

		// мы не должны никогда попадать на цифру, если предыдущий символ не экранируюший
		if unicode.IsDigit(char) && (i == 0 || slice[i-1] != esc) {
			return "", ErrInvalidString
		}

		// если последний символ, прибавляем его и выходим
		if i == last {
			result.WriteRune(char)
			break
		}

		next := slice[i+1]

		switch {
		case next == esc: // следующий символ экранирующий
			result.WriteRune(char)
			i++

		case unicode.IsDigit(next): // следующий символ цифра
			num, err := strconv.Atoi(string(next))
			if err != nil {
				return "", err
			}
			result.WriteString(strings.Repeat(string(char), num))
			i++

		default: // следующий символ не цифра и не экранирующий
			result.WriteRune(char)
		}
	}

	return result.String(), nil
}
