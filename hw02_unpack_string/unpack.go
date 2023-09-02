package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(src string) (string, error) {
	var prev rune
	var bld strings.Builder

	for _, current := range src {
		// цифра
		if unicode.IsDigit(current) {
			// нет предыдущего - ошибка
			if prev == 0 {
				return "", ErrInvalidString
			}

			// есть предыдущий - повторяем n раз
			cnt, err := strconv.Atoi(string(current))
			if err != nil {
				return "", err
			}
			bld.WriteString(strings.Repeat(string(prev), cnt))
			prev = 0
			continue
		}

		// не цифра
		if prev != 0 {
			bld.WriteRune(prev)
		}
		prev = current
	}

	if prev != 0 {
		bld.WriteRune(prev)
	}

	return bld.String(), nil
}
