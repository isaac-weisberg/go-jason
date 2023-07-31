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
			return nil, e("parsing root object failed")
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
			default:
				panic("unhandled parse obj state")
			}
		case jsonWhitespaceTokenType:
			panic("I said, first non-whitespace, please!")
		case jsonColonTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, e("expected a start of a json object, got colon instead")
			case PJOGotKeyWaitingForSeparator:
				// das gud

				// TODOOO NEEDS TO BE IMPLEMENTED NEXT
			default:
				panic("unhandled parse obj state")
			}
		case jsonCurlyOpenBracketTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, e("expected a start of a json object, got curly open bracket. What are you, trying to make an object inside an object?")
			case PJOGotKeyWaitingForSeparator:
				return nil, e("expected colon, but there was suddenly a new object start")
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
			default:
				panic("unhandled parse obj state")
			}
		case jsonCommaTokenType:
			switch state {
			case PJOWaitingForKey:
				return nil, e("expected a start of a json object, got a comma")
			case PJOGotKeyWaitingForSeparator:
				return nil, e("expected colon, but there was suddenly a comma")
			default:
				panic("unhandled parse obj state")
			}
		default:
			panic("token type unhandled")
		}
	}
}
