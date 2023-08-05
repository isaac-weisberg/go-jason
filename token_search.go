package gojason

import (
	"fmt"

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
)

type token struct {
	tokenType tokenType
	payload   []byte
	start     int
	end       int
}

func (token *token) getStartEndString() string {
	return fmt.Sprintf("<%v:%v>", token.start, token.end)
}

func newToken(tokenType tokenType, payload []byte, start int, end int) *token {
	return &token{
		tokenType: tokenType,
		payload:   payload,
		start:     start,
		end:       end,
	}
}

type tokenSearchState int

const (
	invalidoTokenSearchState tokenSearchState = iota
	numberSignStartedTokenSearchState
	numberMaybeTokenSearchState
	whitespaceMaybeTokenSearchState
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
	default:
		panic("RTT unhandled")
	}

	if newState == invalidoTokenSearchState {
		return newFindTokenError(util.E(fmt.Sprintf("failed to interpret byte while parsing token byte = %q", startingByte)))
	}

	state = newState

	var i = start

	var createFindTokenSuccess = func(tokenType tokenType) findTokenResult {
		tokenSearch.byteOffset = i

		var tokenPayload = payload[start:i]
		var token = newToken(tokenType, tokenPayload, start, i)
		return newFindTokenSuccess(token)
	}

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
			case ColonRRT:
				return newFindTokenError(util.E("unexpected colon while we've just gotten a minus"))
			case CurlyOpenBracketRRT:
				return newFindTokenError(util.E("unexpected curly closing bracket while we've just gotten a minus"))
			case CurlyClosingBracketRRT:
				return newFindTokenError(util.E("unexpected curly open bracket while we've just gotten a minus"))
			case CommaRRT:
				return newFindTokenError(util.E("unexpected comma while we've just gotten a minus"))
			default:
				panic("RTT unhandled")
			}
		case numberMaybeTokenSearchState:
			switch byteType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT:
				return createFindTokenSuccess(jsonNumberTokenType)
			case MinusRRT:
				return newFindTokenError(util.E("unexpected minus while the number is already going"))
			case DigitRRT:
				continue
			case ColonRRT:
				return createFindTokenSuccess(jsonNumberTokenType)
			case CurlyOpenBracketRRT:
				return createFindTokenSuccess(jsonNumberTokenType)
			case CurlyClosingBracketRRT:
				return createFindTokenSuccess(jsonNumberTokenType)
			case CommaRRT:
				return createFindTokenSuccess(jsonNumberTokenType)
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
				return createFindTokenSuccess(jsonWhitespaceTokenType)
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
		return createFindTokenSuccess(jsonNumberTokenType)
	case whitespaceMaybeTokenSearchState:
		return createFindTokenSuccess(jsonWhitespaceTokenType)
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
