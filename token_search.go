package gojason

import (
	"fmt"
	"strings"

	"github.com/isaac-weisberg/go-jason/util"
)

type tokenSearch struct {
	payload          []byte
	payloadByteCount int
	byteOffset       int
}

func newTokenSearch(bytes []byte) tokenSearch {
	return tokenSearch{
		payload:          bytes,
		payloadByteCount: len(bytes),
		byteOffset:       0,
	}
}

type tokenType int

const (
	invalidoTokenType tokenType = iota
	jsonNumberTokenType
	jsonWhitespaceTokenType
	jsonColonTokenType
	jsonCurlyOpenBracketTokenType
	jsonCurlyClosingBracketTokenType
	jsonCommaTokenType
	jsonStringTokenType
)

type token struct {
	tokenType   tokenType
	payload     []byte
	start       int
	end         int
	stringValue *string
}

func (token *token) getStartEndString() string {
	return fmt.Sprintf("<%v:%v>", token.start, token.end)
}

func newToken(tokenType tokenType, payload []byte, start int, end int) *token {
	return &token{
		tokenType:   tokenType,
		payload:     payload,
		start:       start,
		end:         end,
		stringValue: nil,
	}
}

func newTokenString(stringValue string, payload []byte, start int, end int) *token {
	return &token{
		tokenType:   jsonStringTokenType,
		payload:     payload,
		start:       start,
		end:         end,
		stringValue: &stringValue,
	}
}

type tokenSearchState int

const (
	invalidoTokenSearchState tokenSearchState = iota
	numberSignStartedTokenSearchState
	numberMaybeTokenSearchState
	whitespaceMaybeTokenSearchState
	stringMaybeTokenSearchState
	stringMaybeButInsideEscapeSequenceSearchState
)

type findTokenResult struct {
	token *token
	err   error
}

func newFindTokenSuccess(token *token) findTokenResult {
	return newFindTokenResult(token, nil)
}

func newFindTokenError(err error) findTokenResult {
	return newFindTokenResult(nil, err)
}

func newFindTokenResult(token *token, err error) findTokenResult {
	return findTokenResult{
		token: token,
		err:   err,
	}
}

func (tokenSearch *tokenSearch) findToken() findTokenResult {
	var start = tokenSearch.byteOffset
	var payloadLen = tokenSearch.payloadByteCount
	var payload = tokenSearch.payload

	if start > payloadLen {
		panic("not supposed to happen")
	}

	if start == payloadLen {
		return newFindTokenSuccess(nil)
	}

	var startingByte byte = payload[start]
	var startingByteType, err = newRecognizedByteType(startingByte)
	if err != nil {
		return newFindTokenError(util.W(err, "failed to deteremine byte type"))
	}

	var state tokenSearchState

	var stringBuilder *strings.Builder

	var newState tokenSearchState
	switch startingByteType {
	case InvalidoRRT:
		panic("how")
	case WhitespaceRRT:
		newState = whitespaceMaybeTokenSearchState
	case MinusRRT:
		newState = numberSignStartedTokenSearchState
	case DigitRRT:
		newState = numberMaybeTokenSearchState
	case ColonRRT:
		var end = start + 1
		tokenSearch.byteOffset = end
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(jsonColonTokenType, tokenPayload, start, end))
	case CurlyOpenBracketRRT:
		var end = start + 1
		tokenSearch.byteOffset = end
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(jsonCurlyOpenBracketTokenType, tokenPayload, start, end))
	case CurlyClosingBracketRRT:
		var end = start + 1
		tokenSearch.byteOffset = end
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(jsonCurlyClosingBracketTokenType, tokenPayload, start, end))
	case CommaRRT:
		var end = start + 1
		tokenSearch.byteOffset = end
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(jsonCommaTokenType, tokenPayload, start, end))
	case DoubleQuoteRRT:
		newState = stringMaybeTokenSearchState
		stringBuilder = &strings.Builder{}
	case BackwardSlashRRT, AnyOtherByteRRT:
		return newFindTokenError(util.E("expected a token start, but got a '%+v'", string(startingByte)))
	default:
		panic("RTT unhandled")
	}

	if newState == invalidoTokenSearchState {
		return newFindTokenError(util.E(fmt.Sprintf("failed to interpret byte while parsing token byte = %q", startingByte)))
	}

	state = newState

	var i = start

	for ; i < payloadLen; i++ {
		var r = payload[i]
		var byteType, err = newRecognizedByteType(r)
		if err != nil {
			return newFindTokenError(util.W(err, "failed to deteremine byte type"))
		}

		switch state {
		case invalidoTokenSearchState:
			panic("how")
		case numberSignStartedTokenSearchState:
			switch byteType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT:
				return newFindTokenError(util.E("unexpected whitespace while we've just gotten a minus"))
			case MinusRRT:
				return newFindTokenError(util.E("unexpected second minus while we've already gotten a minus"))
			case DigitRRT:
				// 'ery nice
				state = numberMaybeTokenSearchState
			case ColonRRT, CurlyOpenBracketRRT, CurlyClosingBracketRRT, CommaRRT, DoubleQuoteRRT, BackwardSlashRRT, AnyOtherByteRRT:
				return newFindTokenError(util.E("a number started with a minus, but we got a '%+v'", string(r)))
			default:
				panic("RTT unhandled")
			}
		case numberMaybeTokenSearchState:
			switch byteType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT, ColonRRT, CurlyOpenBracketRRT, CurlyClosingBracketRRT, CommaRRT:
				tokenSearch.byteOffset = i
				var payloadStart = start
				var payloadEnd = i
				var tokenPayload = payload[payloadStart:payloadEnd]
				var token = newToken(jsonNumberTokenType, tokenPayload, payloadStart, payloadEnd)

				return newFindTokenSuccess(token)
			case MinusRRT:
				return newFindTokenError(util.E("unexpected minus while the number is already going"))
			case DigitRRT:
				continue
			case DoubleQuoteRRT, BackwardSlashRRT, AnyOtherByteRRT:
				return newFindTokenError(util.E("expected digit, but got %+v", string(r)))
			default:
				panic("RTT unhandled")
			}
		case whitespaceMaybeTokenSearchState:
			switch byteType {
			case InvalidoRRT:
				panic("no")
			case WhitespaceRRT:
				continue
			default:
				tokenSearch.byteOffset = i
				var payloadStart = start
				var payloadEnd = i
				var tokenPayload = payload[payloadStart:payloadEnd]
				var token = newToken(jsonWhitespaceTokenType, tokenPayload, payloadStart, payloadEnd)

				return newFindTokenSuccess(token)
			}
		case stringMaybeTokenSearchState:
			switch byteType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT, MinusRRT, DigitRRT, ColonRRT, CurlyOpenBracketRRT, CurlyClosingBracketRRT, CommaRRT, AnyOtherByteRRT:
				stringBuilder.WriteByte(r)
			case DoubleQuoteRRT:
				var resultingString = stringBuilder.String()
				tokenSearch.byteOffset = i + 1
				var payloadStart = start
				var payloadEnd = i
				var tokenPayload = payload[payloadStart:payloadEnd]
				var token = newTokenString(resultingString, tokenPayload, payloadStart, payloadEnd)

				return newFindTokenSuccess(token)
			case BackwardSlashRRT:
				state = stringMaybeButInsideEscapeSequenceSearchState
			default:
				panic("RTT unhandled")
			}
		case stringMaybeButInsideEscapeSequenceSearchState:
			switch byteType {
			case InvalidoRRT:
				panic("nope")
			case WhitespaceRRT, DigitRRT, MinusRRT, ColonRRT, CurlyOpenBracketRRT, CurlyClosingBracketRRT, CommaRRT, AnyOtherByteRRT:
				return newFindTokenError(util.E("string was ungoing with an escape sequence, but got character %+v after backslash", string(r)))
			case DoubleQuoteRRT:
				stringBuilder.WriteByte('"')
			case BackwardSlashRRT:
				stringBuilder.WriteByte('\\')
			}
		default:
			panic("unhandled token search state")
		}
	}

	// Loop finished? that's weird. Must've ran out of payload then
	switch state {
	case invalidoTokenSearchState:
		panic("how")
	case numberSignStartedTokenSearchState:
		return newFindTokenError(util.E("number was started with a sign, but the payload abruptly ended"))
	case numberMaybeTokenSearchState:
		tokenSearch.byteOffset = payloadLen
		var tokenStart = start
		var tokenEnd = payloadLen
		var tokenPayload = payload[tokenStart:tokenEnd]
		var token = newToken(jsonNumberTokenType, tokenPayload, tokenStart, tokenEnd)

		return newFindTokenSuccess(token)
	case whitespaceMaybeTokenSearchState:
		tokenSearch.byteOffset = payloadLen
		var tokenStart = start
		var tokenEnd = payloadLen
		var tokenPayload = payload[tokenStart:tokenEnd]
		var token = newToken(jsonNumberTokenType, tokenPayload, tokenStart, tokenEnd)

		return newFindTokenSuccess(token)
	case stringMaybeTokenSearchState:
		return newFindTokenError(util.E("string was started, but it abruptly ended"))
	case stringMaybeButInsideEscapeSequenceSearchState:
		return newFindTokenError(util.E("string was started and waiting for escape sequence, but it abrubtly ended"))
	default:
		panic("unhandled token search state")
	}
}

func (tokenSearch *tokenSearch) findNonWhitespaceToken() findTokenResult {
	for {
		var findTokenResult = tokenSearch.findToken()

		if findTokenResult.err != nil {
			return findTokenResult
		}

		if findTokenResult.token == nil {
			return findTokenResult
		}

		switch findTokenResult.token.tokenType {
		case jsonWhitespaceTokenType:
			continue
		default:
			return findTokenResult
		}
	}
}
