package gojason

import (
	"github.com/isaac-weisberg/go-jason/util"
	"github.com/isaac-weisberg/go-jason/values"
)

func Parse(jsonString string) (*values.JsonValueAny, error) {
	tokenSearch := newTokenSearch(jsonString)

	var firstTokenSearchResult = tokenSearch.findNonWhitespaceToken()

	var err = firstTokenSearchResult.err
	if err != nil {
		return nil, util.W(err, "token search failed")
	}

	var firstToken = firstTokenSearchResult.token

	if firstToken == nil {
		return nil, util.E("haven't found any token in this jsonString")
	}

	switch firstToken.tokenType {
	case invalidoTokenType:
		panic("what?")
	case jsonNumberTokenType:
		panic("sorry, no number top level objects for now")
	case jsonWhitespaceTokenType:
		panic("I called findNonWhitespaceToken ;)")
	case jsonColonTokenType:
		return nil, util.E("colon can not be the first token in a json")
	case jsonCurlyOpenBracketTokenType:
		// now we're talking
		var parsedObject, err = parseJsonObjectAfterItJustStarted(&tokenSearch)

		if err != nil {
			return nil, util.W(err, "parsing root object failed")
		}

		var any = parsedObject.AsAny()

		return &any, nil
	case jsonCurlyClosingBracketTokenType:
		return nil, util.E("curly closing bracket can not be the first token in a json")
	case jsonCommaTokenType:
		return nil, util.E("comma can not be the first token in a json")
	default:
		panic("unhandled token type")
	}
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=parseJsonObjectState -output parse_json_object_state_string.go
type parseJsonObjectState int

const (
	PJOWaitingForKey parseJsonObjectState = iota
	PJOGotKeyWaitingForSeparator
	PJOGotKeyAndSeparatorWaitingForValue
	PJOGotValueWaitingForCommaOrEnd
	PJOWaitingForKeyAfterComma
)

func parseJsonObjectAfterItJustStarted(tokenSearch *tokenSearch) (*values.JsonValueObject, error) {
	var _ = newParseJsonStateChain()

	var state = PJOWaitingForKey
	var keyValuePairs = values.NewJsonValueObjectKeyValues(0)
	var parsedKey *values.JsonValueAny

	for {
		var tokenSearchResult = tokenSearch.findNonWhitespaceToken()

		var err = tokenSearchResult.err
		if err != nil {
			return nil, util.W(err, "token search failed")
		}

		var token = tokenSearchResult.token
		if token == nil {
			switch state {
			case PJOWaitingForKey:
				return nil, util.E("unexpected lack of token while waiting for first field in the object")
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("unexpected lack of token while waiting for separator after key")
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, util.E("unexpected lack of token while waiting for value")
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, util.E("unexpected lack of token while waiting for end of object")
			case PJOWaitingForKeyAfterComma:
				return nil, util.E("unexpected lack of token while waiting for next key after comma")
			default:
				panic("unhandle parse json object state")
			}
		}

		switch token.tokenType {
		case invalidoTokenType:
			panic("noo")
		case jsonNumberTokenType:
			switch state {
			case PJOWaitingForKey:
				var numberValue = values.NewJsonValueNumber(token.payload)
				var anyValue = numberValue.AsAny()
				parsedKey = &anyValue
				state = PJOGotKeyWaitingForSeparator
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected colon, but there was suddenly a number after a key")
			case PJOGotKeyAndSeparatorWaitingForValue:
				var number = values.NewJsonValueNumber(token.payload)
				var anyValue = number.AsAny()

				var key = parsedKey
				parsedKey = nil

				if key == nil {
					panic("not supposed to happen in this branch")
				}

				keyValuePairs[*key] = anyValue

				state = PJOGotValueWaitingForCommaOrEnd
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, util.E("expected comma or curly closing bracket, but got a number")
			case PJOWaitingForKeyAfterComma:
				var numberValue = values.NewJsonValueNumber(token.payload)
				var anyValue = numberValue.AsAny()
				parsedKey = &anyValue
				state = PJOGotKeyWaitingForSeparator
			default:
				panic("unhandled parse obj state")
			}
		case jsonWhitespaceTokenType:
			panic("I said, first non-whitespace, please!")
		case jsonColonTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, util.E("expected a start of a key, got colon instead at loc=%v", token.getStartEndString())
			case PJOGotKeyWaitingForSeparator:
				// das gud
				state = PJOGotKeyAndSeparatorWaitingForValue
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, util.E("expected value for key %v, but got a colon", parsedKey)
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, util.E("expected comma or curly closing bracket, but got colon")
			case PJOWaitingForKeyAfterComma:
				return nil, util.E("expected next key after comma, but got colon")
			default:
				panic("unhandled parse obj state")
			}
		case jsonCurlyOpenBracketTokenType:
			switch state {
			case PJOWaitingForKey:
				// object as key - interesting

				var object, err = parseJsonObjectAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, util.W(err, "tried to parse key object, but it failed")
				}
				var anyObject = object.AsAny()

				parsedKey = &anyObject
				state = PJOGotKeyWaitingForSeparator
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected colon, but there was suddenly a new object start")
			case PJOGotKeyAndSeparatorWaitingForValue:
				// the value is object then...

				var object, err = parseJsonObjectAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, util.W(err, "tried to parse value object, but it failed")
				}

				var key = parsedKey
				parsedKey = nil

				if key == nil {
					panic("das not supposed to happon")
				}

				var anyValue = object.AsAny()

				keyValuePairs[*key] = anyValue

				state = PJOGotValueWaitingForCommaOrEnd
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, util.E("expected comma or end of object, but got a curly open bracket")
			case PJOWaitingForKeyAfterComma:
				// object as key - interesting

				var object, err = parseJsonObjectAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, util.W(err, "tried to parse key object after comma, parsing failed")
				}
				var anyObject = object.AsAny()

				parsedKey = &anyObject
				state = PJOGotKeyWaitingForSeparator
			default:
				panic("unhandled parse obj state")
			}
		case jsonCurlyClosingBracketTokenType:
			switch state {
			case PJOWaitingForKey:
				// welp
				var jsonObject = values.NewJsonValueObject(keyValuePairs)

				return &jsonObject, nil
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected colon, but there was suddenly an end of the current object")
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, util.E("expected value for key %v, but got curly closing bracket", parsedKey)
			case PJOGotValueWaitingForCommaOrEnd:
				// phew, object ended
				var object = values.NewJsonValueObject(keyValuePairs)
				return &object, nil
			case PJOWaitingForKeyAfterComma:
				return nil, util.E("expected next key after comma, but got curly closing bracket and object is closing")
			default:
				panic("unhandled parse obj state")
			}
		case jsonCommaTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, util.E("expected a start of a json object, got a comma")
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected colon, but there was suddenly a comma")
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, util.E("expected value for key %v, but got a comma", parsedKey)
			case PJOGotValueWaitingForCommaOrEnd:
				// nice, next
				state = PJOWaitingForKeyAfterComma
			case PJOWaitingForKeyAfterComma:
				return nil, util.E("expected next key after comma, but got another comma lol")
			default:
				panic("unhandled parse obj state")
			}
		default:
			panic("token type unhandled")
		}
	}
}
