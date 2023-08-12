package tokenizer

import (
	"fmt"
	"strings"

	"github.com/isaac-weisberg/go-jason/util"
)

type TokenSearch struct {
	Payload          []byte
	PayloadByteCount int
	ByteOffset       int
}

func NewTokenSearch(bytes []byte) TokenSearch {
	return TokenSearch{
		Payload:          bytes,
		PayloadByteCount: len(bytes),
		ByteOffset:       0,
	}
}

type TokenType int

const (
	InvalidoTokenType TokenType = iota
	JsonNumberTokenType
	JsonWhitespaceTokenType
	JsonColonTokenType
	JsonCurlyOpenBracketTokenType
	JsonCurlyClosingBracketTokenType
	JsonCommaTokenType
	JsonStringTokenType
	JsonSquareOpenBracketTokenType
	JsonSquareClosingBracketTokenType
)

type Token struct {
	TokenType   TokenType
	Payload     []byte
	Start       int
	End         int
	StringValue *string
}

func (token *Token) GetStartEndString() string {
	return fmt.Sprintf("<%v:%v>", token.Start, token.End)
}

func newToken(tokenType TokenType, payload []byte, start int, end int) *Token {
	return &Token{
		TokenType:   tokenType,
		Payload:     payload,
		Start:       start,
		End:         end,
		StringValue: nil,
	}
}

func newTokenString(stringValue string, payload []byte, start int, end int) *Token {
	return &Token{
		TokenType:   JsonStringTokenType,
		Payload:     payload,
		Start:       start,
		End:         end,
		StringValue: &stringValue,
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

type FindTokenResult struct {
	Token *Token
	Err   error
}

func newFindTokenSuccess(token *Token) FindTokenResult {
	return newFindTokenResult(token, nil)
}

func newFindTokenError(err error) FindTokenResult {
	return newFindTokenResult(nil, err)
}

func newFindTokenResult(token *Token, err error) FindTokenResult {
	return FindTokenResult{
		Token: token,
		Err:   err,
	}
}

func (tokenSearch *TokenSearch) findToken() FindTokenResult {
	var start = tokenSearch.ByteOffset
	var payloadLen = tokenSearch.PayloadByteCount
	var payload = tokenSearch.Payload

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
		tokenSearch.updateByteOffset(end)
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(JsonColonTokenType, tokenPayload, start, end))
	case CurlyOpenBracketRRT:
		var end = start + 1
		tokenSearch.updateByteOffset(end)
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(JsonCurlyOpenBracketTokenType, tokenPayload, start, end))
	case CurlyClosingBracketRRT:
		var end = start + 1
		tokenSearch.updateByteOffset(end)
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(JsonCurlyClosingBracketTokenType, tokenPayload, start, end))
	case CommaRRT:
		var end = start + 1
		tokenSearch.updateByteOffset(end)
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(JsonCommaTokenType, tokenPayload, start, end))
	case DoubleQuoteRRT:
		newState = stringMaybeTokenSearchState
		stringBuilder = &strings.Builder{}
	case BackwardSlashRRT, AnyOtherByteRRT:
		return newFindTokenError(util.E("expected a token start, but got a '%+v'", string(startingByte)))
	case SquareOpenBracketRRT:
		var end = start + 1
		tokenSearch.updateByteOffset(end)
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(JsonSquareOpenBracketTokenType, tokenPayload, start, end))
	case SquareClosingBracketRRT:
		var end = start + 1
		tokenSearch.updateByteOffset(end)
		var tokenPayload = payload[start:end]
		return newFindTokenSuccess(newToken(JsonSquareClosingBracketTokenType, tokenPayload, start, end))
	default:
		panic("RTT unhandled")
	}

	if newState == invalidoTokenSearchState {
		return newFindTokenError(util.E(fmt.Sprintf("failed to interpret byte while parsing token byte = %q", startingByte)))
	}

	state = newState

	var i = start + 1

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
				tokenSearch.updateByteOffset(i)
				var payloadStart = start
				var payloadEnd = i
				var tokenPayload = payload[payloadStart:payloadEnd]
				var token = newToken(JsonNumberTokenType, tokenPayload, payloadStart, payloadEnd)

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
				tokenSearch.updateByteOffset(i)
				var payloadStart = start
				var payloadEnd = i
				var tokenPayload = payload[payloadStart:payloadEnd]
				var token = newToken(JsonWhitespaceTokenType, tokenPayload, payloadStart, payloadEnd)

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
				tokenSearch.updateByteOffset(i + 1)
				var payloadStart = start + 1
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
				state = stringMaybeTokenSearchState
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
		tokenSearch.updateByteOffset(payloadLen)
		var tokenStart = start
		var tokenEnd = payloadLen
		var tokenPayload = payload[tokenStart:tokenEnd]
		var token = newToken(JsonNumberTokenType, tokenPayload, tokenStart, tokenEnd)

		return newFindTokenSuccess(token)
	case whitespaceMaybeTokenSearchState:
		tokenSearch.updateByteOffset(payloadLen)
		var tokenStart = start
		var tokenEnd = payloadLen
		var tokenPayload = payload[tokenStart:tokenEnd]
		var token = newToken(JsonNumberTokenType, tokenPayload, tokenStart, tokenEnd)

		return newFindTokenSuccess(token)
	case stringMaybeTokenSearchState:
		return newFindTokenError(util.E("string was started, but it abruptly ended"))
	case stringMaybeButInsideEscapeSequenceSearchState:
		return newFindTokenError(util.E("string was started and waiting for escape sequence, but it abrubtly ended"))
	default:
		panic("unhandled token search state")
	}
}

func (tokenSearch *TokenSearch) updateByteOffset(offset int) {
	tokenSearch.ByteOffset = offset
}

func (tokenSearch *TokenSearch) FindNonWhitespaceToken() FindTokenResult {
	for {
		var findTokenResult = tokenSearch.findToken()

		if findTokenResult.Err != nil {
			return findTokenResult
		}

		if findTokenResult.Token == nil {
			return findTokenResult
		}

		switch findTokenResult.Token.TokenType {
		case JsonWhitespaceTokenType:
			continue
		default:
			return findTokenResult
		}
	}
}
