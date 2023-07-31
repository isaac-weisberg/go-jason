package gojason

import "unicode"

type RecognizedRuneType int

const (
	InvalidoRRT RecognizedRuneType = iota
	WhitespaceRRT
	DigitRRT
	MinusRRT
	ColonRRT
	CurlyOpenBracketRRT
	CurlyClosingBracketRRT
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

	if r == '{' {
		return CurlyOpenBracketRRT, nil
	}

	if r == '}' {
		return CurlyClosingBracketRRT, nil
	}

	return InvalidoRRT, e("rune type unrecognized %q", r)
}
