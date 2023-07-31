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
	jsonCommaTokenType
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

	if start < payloadLen {
		return newFindTokenSuccess(nil)
	}

	var startRune rune = payload[start]
	var runeType, err = newRuneType(startRune)
	if err != nil {
		return newFindTokenError(w(err, "failed to deteremine rune type"))
	}

	var state tokenSearchState

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
		tokenSearch.runeOffset = start + 1
		var tokenPayload = payload[start : start+1]
		return newFindTokenSuccess(newToken(jsonColonTokenType, tokenPayload))
	case CurlyOpenBracketRRT:
		tokenSearch.runeOffset = start + 1
		var tokenPayload = payload[start : start+1]
		return newFindTokenSuccess(newToken(jsonCurlyOpenBracketTokenType, tokenPayload))
	case CurlyClosingBracketRRT:
		tokenSearch.runeOffset = start + 1
		var tokenPayload = payload[start : start+1]
		return newFindTokenSuccess(newToken(jsonCurlyClosingBracketTokenType, tokenPayload))
	case CommaRRT:
		tokenSearch.runeOffset = start + 1
		var tokenPayload = payload[start : start+1]
		return newFindTokenSuccess(newToken(jsonCommaTokenType, tokenPayload))
	default:
		panic("RTT unhandled")
	}

	if newState == invalidoTokenSearchState {
		return newFindTokenError(e(fmt.Sprintf("failed to interpret rune while parsing token rune = %q", startRune)))
	}

	state = newState

	var i = start

	var createFindTokenSuccess = func(tokenType tokenType) findTokenResult {
		tokenSearch.runeOffset = i

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
			case CommaRRT:
				return newFindTokenError(e("unexpected comma while we've just gotten a minus"))
			default:
				panic("RTT unhandled")
			}
		case numberMaybeTokenSearchState:
			switch runeType {
			case InvalidoRRT:
				panic("how")
			case WhitespaceRRT:
				return createFindTokenSuccess(jsonNumberTokenType)
			case MinusRRT:
				return newFindTokenError(e("unexpected minus while the number is already going"))
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
			switch runeType {
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
		return newFindTokenError(e("number was started with a sign, but the payload abruptly ended"))
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
