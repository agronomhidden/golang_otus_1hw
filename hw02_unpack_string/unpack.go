package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"

	"unicode/utf8"
)

const slash = `\`

var ErrInvalidString = errors.New("invalid string")
var ErrUnexpectedEnd = errors.New("unexpected end of string")

type symbolCollection struct {
	str string
}

func (e *symbolCollection) IsEmpty() bool {
	return e.str == ""
}
func (e *symbolCollection) ExtractSymbol() string {
	symbol, runeLength := utf8.DecodeRuneInString(e.str)
	e.str = e.str[runeLength:]
	return string(symbol)
}

func Unpack(input string) (string, error) {
	var collection = symbolCollection{str: input}
	var result strings.Builder
	var buffer = ""

	for !collection.IsEmpty() {
		symbol := collection.ExtractSymbol()
		number, err := strconv.Atoi(symbol)
		if symbol == slash {
			if collection.IsEmpty() {
				return "", ErrUnexpectedEnd
			}
			symbol = collection.ExtractSymbol()
		}
		if err != nil {
			result.WriteString(buffer)
			buffer = symbol
			continue
		}
		if number == 0 || buffer == "" {
			return "", ErrInvalidString
		}
		result.WriteString(strings.Repeat(buffer, number))
		buffer = ""
	}
	result.WriteString(buffer)
	return result.String(), nil
}
