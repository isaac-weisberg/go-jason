package gojason

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
	BackwardSlashRRT
	AnyOtherByteRRT
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

	if r == '\\' {
		return BackwardSlashRRT, nil
	}

	return AnyOtherByteRRT, nil
}
