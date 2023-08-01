package gojason

func parse(jsonString string) (*JsonValueAny, error) {
	tokenSearch := newTokenSearch(jsonString)

	var firstTokenSearchResult = tokenSearch.findNonWhitespaceToken()

	var err = firstTokenSearchResult.err
	if err != nil {
		return nil, w(err, "token search failed")
	}

	var firstToken = firstTokenSearchResult.token

	if firstToken == nil {
		return nil, e("haven't found any token in this jsonString")
	}

	switch firstToken.tokenType {
	case invalidoTokenType:
		panic("what?")
	case jsonNumberTokenType:
		panic("sorry, no number top level objects for now")
	case jsonWhitespaceTokenType:
		panic("I called findNonWhitespaceToken ;)")
	case jsonColonTokenType:
		return nil, e("colon can not be the first token in a json")
	case jsonCurlyOpenBracketTokenType:
		// now we're talking
		var parsedObject, err = parseJsonObjectAfterItJustStarted(&tokenSearch)

		if err != nil {
			return nil, w(err, "parsing root object failed")
		}

		var any = parsedObject.asAny()

		return &any, nil
	case jsonCurlyClosingBracketTokenType:
		return nil, e("curly closing bracket can not be the first token in a json")
	case jsonCommaTokenType:
		return nil, e("comma can not be the first token in a json")
	default:
		panic("unhandled token type")
	}
}

type parseJsonObjectState int

const (
	PJOWaitingForKey parseJsonObjectState = iota
	PJOGotKeyWaitingForSeparator
	PJOGotKeyAndSeparatorWaitingForValue
	PJOGotValueWaitingForCommaOrEnd
	PJOWaitingForKeyAfterComma
)

func parseJsonObjectAfterItJustStarted(tokenSearch *tokenSearch) (*JsonValueObject, error) {
	for {
		var tokenSearchResult = tokenSearch.findNonWhitespaceToken()

		var err = tokenSearchResult.err
		if err != nil {
			return nil, w(err, "token search failed")
		}

		var state = PJOWaitingForKey
		var keyValuePairs = newJsonValueObjectKeyValues()
		var parsedKey *JsonValueAny

		var token = tokenSearchResult.token
		if token == nil {
			switch state {
			case PJOWaitingForKey:
				return nil, e("unexpected lack of token while waiting for first field in the object")
			case PJOGotKeyWaitingForSeparator:
				return nil, e("unexpected lack of token while waiting for separator after key")
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, e("unexpected lack of token while waiting for value")
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, e("unexpected lack of token while waiting for end of object")
			case PJOWaitingForKeyAfterComma:
				return nil, e("unexpected lack of token while waiting for next key after comma")
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
				var numberValue = newJsonValueNumber(token.payload)
				var anyValue = numberValue.asAny()
				parsedKey = &anyValue
				state = PJOGotKeyWaitingForSeparator
			case PJOGotKeyWaitingForSeparator:
				return nil, e("expected colon, but there was suddenly a number after a key")
			case PJOGotKeyAndSeparatorWaitingForValue:
				var number = newJsonValueNumber(token.payload)
				var anyValue = number.asAny()

				var key = parsedKey
				parsedKey = nil

				if key == nil {
					panic("not supposed to happen in this branch")
				}

				keyValuePairs[*key] = anyValue

				state = PJOGotValueWaitingForCommaOrEnd
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, e("expected comma or curly closing bracket, but got a number")
			case PJOWaitingForKeyAfterComma:
				var numberValue = newJsonValueNumber(token.payload)
				var anyValue = numberValue.asAny()
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
				return nil, e("expected a start of a key, got colon instead")
			case PJOGotKeyWaitingForSeparator:
				// das gud
				state = PJOGotKeyAndSeparatorWaitingForValue
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, e("expected value for key %v, but got a colon", parsedKey)
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, e("expected comma or curly closing bracket, but got colon")
			case PJOWaitingForKeyAfterComma:
				return nil, e("expected next key after comma, but got colon")
			default:
				panic("unhandled parse obj state")
			}
		case jsonCurlyOpenBracketTokenType:
			switch state {
			case PJOWaitingForKey:
				// object as key - interesting

				var object, err = parseJsonObjectAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, w(err, "tried to parse key object, but it failed")
				}
				var anyObject = object.asAny()

				parsedKey = &anyObject
				state = PJOGotKeyWaitingForSeparator
			case PJOGotKeyWaitingForSeparator:
				return nil, e("expected colon, but there was suddenly a new object start")
			case PJOGotKeyAndSeparatorWaitingForValue:
				// the value is object then...

				var object, err = parseJsonObjectAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, w(err, "tried to parse value object, but it failed")
				}

				var key = parsedKey
				parsedKey = nil

				if key == nil {
					panic("das not supposed to happon")
				}

				var anyValue = object.asAny()

				keyValuePairs[*key] = anyValue

				state = PJOGotValueWaitingForCommaOrEnd
			case PJOGotValueWaitingForCommaOrEnd:
				return nil, e("expected comma or end of object, but got a curly open bracket")
			case PJOWaitingForKeyAfterComma:
				// object as key - interesting

				var object, err = parseJsonObjectAfterItJustStarted(tokenSearch)
				if err != nil {
					return nil, w(err, "tried to parse key object after comma, parsing failed")
				}
				var anyObject = object.asAny()

				parsedKey = &anyObject
				state = PJOGotKeyWaitingForSeparator
			default:
				panic("unhandled parse obj state")
			}
		case jsonCurlyClosingBracketTokenType:
			switch state {
			case PJOWaitingForKey:
				// welp
				var jsonObject = newJsonValueObject(keyValuePairs)

				return &jsonObject, nil
			case PJOGotKeyWaitingForSeparator:
				return nil, e("expected colon, but there was suddenly an end of the current object")
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, e("expected value for key %v, but got curly closing bracket", parsedKey)
			case PJOGotValueWaitingForCommaOrEnd:
				// phew, object ended
				var object = newJsonValueObject(keyValuePairs)
				return &object, nil
			case PJOWaitingForKeyAfterComma:
				return nil, e("expected next key after comma, but got curly closing bracket and object is closing")
			default:
				panic("unhandled parse obj state")
			}
		case jsonCommaTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, e("expected a start of a json object, got a comma")
			case PJOGotKeyWaitingForSeparator:
				return nil, e("expected colon, but there was suddenly a comma")
			case PJOGotKeyAndSeparatorWaitingForValue:
				return nil, e("expected value for key %v, but got a comma", parsedKey)
			case PJOGotValueWaitingForCommaOrEnd:
				// nice, next
				state = PJOWaitingForKeyAfterComma
			case PJOWaitingForKeyAfterComma:
				return nil, e("expected next key after comma, but got another comma lol")
			default:
				panic("unhandled parse obj state")
			}
		default:
			panic("token type unhandled")
		}
	}
}
