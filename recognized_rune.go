package gojason

import (
	"github.com/isaac-weisberg/go-jason/util"
)

type RecognizedByteType int

const (
	InvalidoRRT RecognizedByteType = iota // RRT because it was previous called RecognizedRuneType
	WhitespaceRRT
	DigitRRT
	MinusRRT
	ColonRRT
	CurlyOpenBracketRRT
	CurlyClosingBracketRRT
	CommaRRT
	DoubleQuoteRRT
)

func isEndOfLine(r byte) bool {
	return r == '\n'
}

func isWhitespace(r byte) bool {
	switch r {
	case '\n':
		return true
	case '\t':
		return true
	case ' ':
		return true
	case '\r':
		return true
	default:
		return false
	}
}

func isDigit(r byte) bool {
	if '0' <= r && r <= '9' {
		return true
	}
	return false
}

func newRecognizedByteType(r byte) (RecognizedByteType, error) {
	if isWhitespace(r) {
		return WhitespaceRRT, nil
	}

	if r == '-' {
		return MinusRRT, nil
	}

	if isDigit(r) {
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

	if r == ',' {
		return CommaRRT, nil
	}

	if r == '"' {
		return DoubleQuoteRRT, nil
	}

	// Add backslash and all string escape sequences

	return InvalidoRRT, util.E("byte type unrecognized %q", r)
}
