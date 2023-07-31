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
				newState = numberMaybeTokenSearchState
			case DigitRRT:
				newState = numberMaybeTokenSearchState
			case ColonRRT:
				return createFindTokenSuccess(jsonColonTokenType)
			case CurlyOpenBracketRRT:
				return createFindTokenSuccess(jsonCurlyOpenBracketTokenType)
			case CurlyClosingBracketRRT:
				return createFindTokenSuccess(jsonCurlyClosingBracketTokenType)
			default:
				panic("RTT unhandled")
			}

			if newState == invalidoTokenSearchState {
				return newFindTokenError(e(fmt.Sprintf("failed to interpret rune while parsing token rune = %q", r)))
			}

			state = newState
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
	default:
		panic("unhandled token search state")
	}
}

type JsonNumberToken struct {
	payload string
}
