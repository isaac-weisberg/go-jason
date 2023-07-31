package gojason

type jsonContextType int

const (
	JCTObjectOpened jsonContextType = iota
)

type jsonContext struct {
	jsonContextType jsonContextType
}

func newJsonContext(jsonContextType jsonContextType) jsonContext {

}

func parse(jsonString string) error {
	tokenSearch := newTokenSearch(jsonString)

	var firstTokenSearchResult = tokenSearch.findNonWhitespaceToken()

	var err = firstTokenSearchResult.err
	if err != nil {
		return w(err, "token search failed")
	}

	var firstToken = firstTokenSearchResult.token

	if firstToken == nil {
		return e("haven't found any token in this jsonString")
	}

	var context []jsonContext

	var 

	switch firstToken.tokenType {
	case invalidoTokenType:
		panic("what?")
	case jsonNumberTokenType:
		panic("sorry, no number top level objects for now")
	case jsonWhitespaceTokenType:
		panic("I called findNonWhitespaceToken ;)")
	case jsonColonTokenType:
		return e("colon can not be the first token in a json")
	case jsonCurlyOpenBracketTokenType:
		// now we're talking
		context = []jsonContext{jsonContext{jsonContextType: JCTObjectOpened}}
	case jsonCurlyClosingBracketTokenType:
		return e("curly closing bracket can not be the first token in a json")
	case jsonCommaTokenType:
		return e("comma can not be the first token in a json")
	default:
		panic("unhandled token type")
	}

	for {
		tokenSearchResult := tokenSearch.findNonWhitespaceToken()

		if tokenSearchResult.err != nil {
			return w(err, "token search failed")
		}

		if tokenSearchResult.token == nil {
			// this means json ended
		}
	}
}
