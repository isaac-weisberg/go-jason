package parser

import (
	"github.com/isaac-weisberg/go-jason/tokenizer"
	"github.com/isaac-weisberg/go-jason/util"
	"github.com/isaac-weisberg/go-jason/values"
)

func Parse(bytes []byte) (*values.JsonValueAny, error) {
	tokenSearch := tokenizer.NewTokenSearch(bytes)

	var firstTokenSearchResult = tokenSearch.FindNonWhitespaceToken()

	var err = firstTokenSearchResult.Err
	if err != nil {
		return nil, util.W(err, "token search failed")
	}

	var firstToken = firstTokenSearchResult.Token

	if firstToken == nil {
		return nil, util.E("haven't found any token in this jsonString")
	}

	switch firstToken.TokenType {
	case tokenizer.InvalidoTokenType:
		panic("what?")
	case tokenizer.JsonNumberTokenType:
		numberValue := values.NewJsonValueNumber(firstToken.Payload)
		any := numberValue.AsAny()
		return &any, nil
	case tokenizer.JsonWhitespaceTokenType:
		panic("I called findNonWhitespaceToken ;)")
	case tokenizer.JsonColonTokenType:
		return nil, util.E("colon can not be the first token in a json")
	case tokenizer.JsonCurlyOpenBracketTokenType:
		// now we're talking
		var parsedObject, err = parseJsonObjectAfterItJustStarted(&tokenSearch)

		if err != nil {
			return nil, util.W(err, "parsing root object failed")
		}

		var any = parsedObject.AsAny()

		return &any, nil
	case tokenizer.JsonCurlyClosingBracketTokenType:
		return nil, util.E("curly closing bracket can not be the first token in a json")
	case tokenizer.JsonCommaTokenType:
		return nil, util.E("comma can not be the first token in a json")
	case tokenizer.JsonStringTokenType:
		stringValue := values.NewJsonValueString(*firstToken.StringValue)
		var any = stringValue.AsAny()
		return &any, nil
	case tokenizer.JsonSquareOpenBracketTokenType:
		return nil, util.E("root arrays are not supported at the time")
	case tokenizer.JsonSquareClosingBracketTokenType:
		return nil, util.E("square closing bracket can not be the first token in a json")
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

func parseJsonObjectAfterItJustStarted(tokenSearch *tokenizer.TokenSearch) (*values.JsonValueObject, error) {
	var _ = newParseJsonStateChain()

	var state = PJOWaitingForKey
	var keyValuePairs = values.NewJsonValueObjectKeyValues(0)
	var parsedKey *values.JsonValueAny

	for {
		var tokenSearchResult = tokenSearch.FindNonWhitespaceToken()

		var err = tokenSearchResult.Err
		if err != nil {
			return nil, util.W(err, "token search failed")
		}

		var token = tokenSearchResult.Token
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

		switch token.TokenType {
		case tokenizer.InvalidoTokenType:
			panic("noo")
		case tokenizer.JsonNumberTokenType:
			switch state {
			case PJOWaitingForKey:
				var numberValue = values.NewJsonValueNumber(token.Payload)
				var anyValue = numberValue.AsAny()
				parsedKey = &anyValue
				state = PJOGotKeyWaitingForSeparator
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected colon, but there was suddenly a number after a key")
			case PJOGotKeyAndSeparatorWaitingForValue:
				var number = values.NewJsonValueNumber(token.Payload)
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
				var numberValue = values.NewJsonValueNumber(token.Payload)
				var anyValue = numberValue.AsAny()
				parsedKey = &anyValue
				state = PJOGotKeyWaitingForSeparator
			default:
				panic("unhandled parse obj state")
			}
		case tokenizer.JsonWhitespaceTokenType:
			panic("I said, first non-whitespace, please!")
		case tokenizer.JsonColonTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, util.E("expected a start of a key, got colon instead at loc=%v", token.GetStartEndString())
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
		case tokenizer.JsonCurlyOpenBracketTokenType:
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
		case tokenizer.JsonCurlyClosingBracketTokenType:
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
		case tokenizer.JsonCommaTokenType:
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
		case tokenizer.JsonStringTokenType:
			switch state {
			case PJOWaitingForKey:
				var stringValue = values.NewJsonValueString(*token.StringValue)
				var anyValue = stringValue.AsAny()
				parsedKey = &anyValue
				state = PJOGotKeyWaitingForSeparator
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected colon, but there was suddenly a number after a key")
			case PJOGotKeyAndSeparatorWaitingForValue:
				var stringValue = values.NewJsonValueString(*token.StringValue)
				var anyValue = stringValue.AsAny()

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
				var stringValue = values.NewJsonValueString(*token.StringValue)
				var anyValue = stringValue.AsAny()
				parsedKey = &anyValue
				state = PJOGotKeyWaitingForSeparator
			default:
				panic("unhandled parse obj state")
			}
		case tokenizer.JsonSquareOpenBracketTokenType:
			switch state {
			case PJOWaitingForKey, PJOWaitingForKeyAfterComma:
				return nil, util.E("expected start of an key for an object, but got a square open bracket, which means, an array")
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected a separator, but got a square open bracket")
			case PJOGotKeyAndSeparatorWaitingForValue:
				// array starting

				arrayValue, err := parseArrayAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, util.W(err, "failed to parse an array value for key")
				}
				var anyValue = arrayValue.AsAny()

				var key = parsedKey
				parsedKey = nil
				if key == nil {
					panic("not supposed to happen in this branch")
				}

				keyValuePairs[*key] = anyValue

				state = PJOGotValueWaitingForCommaOrEnd
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, util.E("expected a comma or end of object, but got a square open bracket")
			default:
				panic("unhandled parse obj state")
			}
		case tokenizer.JsonSquareClosingBracketTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, util.E("expected an object key, but got a square closing bracket")
			case PJOGotKeyWaitingForSeparator:
				return nil, util.E("expected separator after key in an object, but got a square closing bracket")
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, util.E("expected value after key and separator in an object, but got a square closing bracket")
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, util.E("expected comma or end of an object, but got a square closing bracket")
			case PJOWaitingForKeyAfterComma:
				return nil, util.E("expected a new key in an object, but got a square closing bracket")
			default:
				panic("unhandled parse obj state")
			}
		default:
			panic("token type unhandled")
		}
	}
}

type ParseArrayState int64

const (
	PAInvalido ParseArrayState = iota
	PAWaitingForFirstValueOrEnd
	PAWaitingForNextValue
	PAWaitingSeparatorOrEnd
)

func parseArrayAfterItJustStarted(tokenSearch *tokenizer.TokenSearch) (*values.JsonValueArray, error) {
	var arrayValueBuffer = make([]values.JsonValueAny, 0)
	var state = PAWaitingForFirstValueOrEnd

	for {
		var tokenRes = tokenSearch.FindNonWhitespaceToken()
		err := tokenRes.Err
		if err != nil {
			return nil, util.W(err, "failed to parse a non whitespace token while parsing array")
		}

		if tokenRes.Token == nil {
			return nil, util.E("was parsing an array, but it abruptly ended...")
		}

		var token = tokenRes.Token
		switch token.TokenType {
		case tokenizer.JsonNumberTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd, PAWaitingForNextValue:
				var numberValue = values.NewJsonValueNumber(token.Payload)
				var anyValue = numberValue.AsAny()

				arrayValueBuffer = append(arrayValueBuffer, anyValue)
				state = PAWaitingSeparatorOrEnd
			case PAWaitingSeparatorOrEnd:
				return nil, util.E("expected an array separator or an array end, got a number")
			default:
				panic("unhandled")
			}
		case tokenizer.JsonWhitespaceTokenType:
			panic("impossibru")
		case tokenizer.JsonColonTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd, PAWaitingForNextValue, PAWaitingSeparatorOrEnd:
				return nil, util.E("was parsing an array but suddenly, got a colon")
			default:
				panic("unhandled")
			}
		case tokenizer.JsonCurlyOpenBracketTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd, PAWaitingForNextValue:
				parsedObject, err := parseJsonObjectAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, util.W(err, "parsing an object for an array value failed")
				}
				var anyValue = parsedObject.AsAny()

				arrayValueBuffer = append(arrayValueBuffer, anyValue)
				state = PAWaitingSeparatorOrEnd
			case PAWaitingSeparatorOrEnd:
				return nil, util.E("expected an array separator or end, but got a curly open bracket")
			default:
				panic("unhandled")
			}
		case tokenizer.JsonCurlyClosingBracketTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd, PAWaitingForNextValue, PAWaitingSeparatorOrEnd:
				return nil, util.E("was parsing an array but suddenly, got a curly closing bracket")
			default:
				panic("unhandled")
			}
		case tokenizer.JsonCommaTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd:
				return nil, util.E("expected a first value of array or an end, but got a comma")
			case PAWaitingForNextValue:
				return nil, util.E("expected a next value of array, but got a comma")
			case PAWaitingSeparatorOrEnd:
				state = PAWaitingForNextValue
			default:
				panic("unhandled")
			}
		case tokenizer.JsonStringTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd, PAWaitingForNextValue:
				var jsonStringValue = values.NewJsonValueString(*token.StringValue)
				var anyValue = jsonStringValue.AsAny()

				arrayValueBuffer = append(arrayValueBuffer, anyValue)
				state = PAWaitingSeparatorOrEnd
			case PAWaitingSeparatorOrEnd:
				return nil, util.E("expected an array separator or array end, but got a string")
			default:
				panic("unhandled")
			}
		case tokenizer.JsonSquareOpenBracketTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd, PAWaitingForNextValue:
				// array o'clock

				arrayValue, err := parseArrayAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, util.W(err, "parsing array value inside the array failed")
				}
				var anyValue = arrayValue.AsAny()

				arrayValueBuffer = append(arrayValueBuffer, anyValue)
				state = PAWaitingSeparatorOrEnd

			case PAWaitingSeparatorOrEnd:
				return nil, util.E("expected an array separator or end, but got a start of new array, square open bracket")
			default:
				panic("unhandled")
			}
		case tokenizer.JsonSquareClosingBracketTokenType:
			switch state {
			case PAWaitingForFirstValueOrEnd, PAWaitingSeparatorOrEnd:
				// array is done, hoooray

				var arrayObject = values.NewJsonValueArray(arrayValueBuffer)
				return &arrayObject, nil
			case PAWaitingForNextValue:
				return nil, util.E("Expected a next object in the array, but got array end, square closing bracket")
			default:
				panic("unhandled")
			}
		default:
			panic("unhandled")
		}
	}

}
