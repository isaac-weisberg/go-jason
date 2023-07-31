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

	var createToken = func(tokenType tokenType) *token {
		var tokenPayload = payload[start:i]
		return newToken(tokenType, tokenPayload)
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
				return newFindTokenSuccess(createToken(jsonColonTokenType))
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
				return newFindTokenSuccess(createToken(jsonNumberTokenType))
			case MinusRRT:
				return newFindTokenError(e("unexpected minus while the number is already going"))
			case DigitRRT:
				continue
			case ColonRRT:
				return newFindTokenError(e("unexpected colon while the number is being parsed"))
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
				return newFindTokenSuccess(createToken(jsonWhitespaceTokenType))
			}
		default:
			panic("unhandled token search state")
		}
	}

	// Loop returned? that's weird. Must've ran out of payload then
	tokenSearch.runeOffset = i

	switch state {
	case invalidoTokenSearchState:
		panic("how")
	case initialTokenSearchState:
		return newFindTokenSuccess(nil)
	case numberMaybeTokenSearchState:
		return newFindTokenSuccess(createToken(jsonNumberTokenType))
	case whitespaceMaybeTokenSearchState:
		return newFindTokenSuccess(createToken(jsonWhitespaceTokenType))
	default:
		panic("unhandled token search state")
	}
}

type JsonNumberToken struct {
	payload string
}
