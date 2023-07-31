package gojason

import "unicode"

type RecognizedRuneType int

const (
	InvalidoRRT RecognizedRuneType = iota
	WhitespaceRRT
	DigitRRT
	MinusRRT
	ColonRRT
)

func isEndOfLine(r rune) bool {
	return r == '\n'
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r) || isEndOfLine(r)
}

func newRuneType(r rune) (RecognizedRuneType, error) {
	if isWhitespace(r) {
		return WhitespaceRRT, nil
	}

	if r == '-' {
		return MinusRRT, nil
	}

	if unicode.IsDigit(r) {
		return DigitRRT, nil
	}

	if r == ':' {
		return ColonRRT, nil
	}

	return InvalidoRRT, e("rune type unrecognized %q", r)
}
