package gojason

import (
	"fmt"
)

type tokenSearch struct {
	payload          []rune
	payloadRuneCount int
	runeOffset       int
}

func newTokenSearch(payload string) tokenSearch {
	var runes = []rune(payload)
	return tokenSearch{
		payload:          runes,
		payloadRuneCount: len(runes),
		runeOffset:       0,
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
)

type token struct {
	tokenType tokenType
	payload   []rune
}

func newToken(tokenType tokenType, payload []rune) *token {
	return &token{
		tokenType: tokenType,
		payload:   payload,
	}
}

type tokenSearchState int

const (
	invalidoTokenSearchState tokenSearchState = iota
	initialTokenSearchState
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
	var start = tokenSearch.runeOffset
	var payloadLen = tokenSearch.payloadRuneCount
	var payload = tokenSearch.payload

	var state tokenSearchState = initialTokenSearchState
	var i = start

	var createFindTokenSuccess = func(tokenType tokenType) findTokenResult {
		var tokenPayload = payload[start:i]
		var token = newToken(tokenType, tokenPayload)
		return newFindTokenSuccess(token)
	}

	for ; i < payloadLen; i++ {
		var r = payload[i]
		var runeType, err = newRuneType(r)
		if err != nil {
			return newFindTokenError(w(err, "failed to deteremine rune type"))
		}

		switch state {
		case invalidoTokenSearchState:
			panic("how")
		case initialTokenSearchState:
			var newState tokenSearchState
			switch runeType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT:
				newState = whitespaceMaybeTokenSearchState
			case MinusRRT:
				newState = numberSignStartedTokenSearchState
			case DigitRRT:
				newState = numberMaybeTokenSearchState
			case ColonRRT:
				tokenSearch.runeOffset = i + 1
				return createFindTokenSuccess(jsonColonTokenType)
			case CurlyOpenBracketRRT:
				tokenSearch.runeOffset = i + 1
				return createFindTokenSuccess(jsonCurlyOpenBracketTokenType)
			case CurlyClosingBracketRRT:
				tokenSearch.runeOffset = i + 1
				return createFindTokenSuccess(jsonCurlyClosingBracketTokenType)
			default:
				panic("RTT unhandled")
			}

			if newState == invalidoTokenSearchState {
				return newFindTokenError(e(fmt.Sprintf("failed to interpret rune while parsing token rune = %q", r)))
			}

			state = newState
		case numberSignStartedTokenSearchState:
			switch runeType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT:
				return newFindTokenError(e("unexpected whitespace while we've just gotten a minus"))
			case MinusRRT:
				return newFindTokenError(e("unexpected second minus while we've already gotten a minus"))
			case DigitRRT:
				// 'ery nice
				state = numberMaybeTokenSearchState
			case ColonRRT:
				return newFindTokenError(e("unexpected colon while we've just gotten a minus"))
			case CurlyOpenBracketRRT:
				return newFindTokenError(e("unexpected curly closing bracket while we've just gotten a minus"))
			case CurlyClosingBracketRRT:
				return newFindTokenError(e("unexpected curly open bracket while we've just gotten a minus"))
			default:
				panic("RTT unhandled")
			}
		case numberMaybeTokenSearchState:
			switch runeType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT:
				tokenSearch.runeOffset = i
				return createFindTokenSuccess(jsonNumberTokenType)
			case MinusRRT:
				return newFindTokenError(e("unexpected minus while the number is already going"))
			case DigitRRT:
				continue
			case ColonRRT:
				tokenSearch.runeOffset = i
				return createFindTokenSuccess(jsonNumberTokenType)
			case CurlyOpenBracketRRT:
				tokenSearch.runeOffset = i
				return createFindTokenSuccess(jsonNumberTokenType)
			case CurlyClosingBracketRRT:
				tokenSearch.runeOffset = i
				return createFindTokenSuccess(jsonNumberTokenType)
			default:
				panic("RTT unhandled")
			}
		case whitespaceMaybeTokenSearchState:
			switch runeType {
			case InvalidoRRT:
				panic("no")
			case WhitespaceRRT:
				continue
			default:
				tokenSearch.runeOffset = i
				return createFindTokenSuccess(jsonWhitespaceTokenType)
			}
		default:
			panic("unhandled token search state")
		}
	}

	// Loop finished? that's weird. Must've ran out of payload then
	tokenSearch.runeOffset = i

	switch state {
	case invalidoTokenSearchState:
		panic("how")
	case initialTokenSearchState:
		return newFindTokenSuccess(nil)
	case numberMaybeTokenSearchState:
		return createFindTokenSuccess(jsonNumberTokenType)
	case whitespaceMaybeTokenSearchState:
		return createFindTokenSuccess(jsonWhitespaceTokenType)
	case numberSignStartedTokenSearchState:
		return newFindTokenError(e("number was started with a sign, but the payload abruptly ended"))
	default:
		panic("unhandled token search state")
	}
}

type JsonNumberToken struct {
	payload string
}
