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
		if unicode.IsDigit(current) {
			// цифра
			if prev != 0 {
				// есть предыдущий - повторяем n раз
				cnt, err := strconv.Atoi(string(current))
				if err != nil {
					return "", err
				}
				bld.WriteString(strings.Repeat(string(prev), cnt))
				prev = 0
			} else {
				// нет предыдущего - ошибка
				return "", ErrInvalidString
			}
		} else {
			// не цифра
			if prev != 0 {
				bld.WriteRune(prev)
			}
			prev = current
		}
	}

	if prev != 0 {
		bld.WriteRune(prev)
	}

	return bld.String(), nil
}
