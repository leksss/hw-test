package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	if inputString == "" {
		return "", nil
	}

	var sb strings.Builder
	pr := int32(0)
	for pos, cr := range inputString {
		if pos == 0 && unicode.IsDigit(cr) {
			return "", ErrInvalidString
		}
		if unicode.IsDigit(cr) && unicode.IsDigit(pr) {
			return "", ErrInvalidString
		}

		repeatCnt := -1
		if pr != 0 {
			if unicode.IsDigit(cr) {
				repeatCnt, _ = strconv.Atoi(string(cr))
			}

			if repeatCnt >= 0 {
				sb.WriteString(strings.Repeat(string(pr), repeatCnt))
			} else if !unicode.IsDigit(pr) {
				sb.WriteRune(pr)
			}
		}
		pr = cr
	}

	lastRune, _ := utf8.DecodeLastRuneInString(inputString)
	if !unicode.IsDigit(lastRune) {
		sb.WriteRune(lastRune)
	}

	return sb.String(), nil
}
