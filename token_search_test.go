package gojason

//go:generate go run golang.org/x/tools/cmd/stringer -type=tokenType -output tokentype_string_test.go

import (
	"fmt"
	"strings"
	"testing"

	"github.com/isaac-weisberg/go-jason/util"
)

func TestSimpleTokenSearch(t *testing.T) {
	jsonString := `
	{
		{
			35: 60
		}: {
			90: 35
		}
	}
	`

	var tokenSearch = newTokenSearch(jsonString)

	var allTokens, err = tokenSearch.findAllTokens()
	if err != nil {
		t.Errorf("find all tokens failed with error = %v", err.Error())
	}

	var allTokensJoined = stringForSlice(allTokens, "")

	if allTokensJoined != jsonString {
		t.Errorf("the resulting tokens sequence is not the same as the source json")
	}

	// fmt.Printf("find all tokens allTokens = \n%v\n", stringForSlice(allTokens, ""))

	// var allTokenTypes = sliceMap[token, string](allTokens, func(t token) string {
	// 	return t.tokenType.String()
	// })
	// var allTokenTypesJoined = strings.Join(allTokenTypes, ", ")
	// fmt.Printf("all token types: \n%v\n", allTokenTypesJoined)
}

func (tokenSearch *tokenSearch) findAllTokens() ([]token, error) {
	var tokens = make([]token, 0)

	for {
		var findTokenResult = tokenSearch.findToken()

		var err = findTokenResult.err
		if err != nil {
			return nil, util.W(err, "find token failed")
		}

		var token = findTokenResult.token

		if token == nil {
			return tokens, nil
		}

		tokens = append(tokens, *token)
	}
}

func (token token) String() string {
	return fmt.Sprintf("%v", string(token.payload))
}

func sliceMap[I interface{}, O interface{}](arr []I, transform func(I) O) []O {
	var newSlice = make([]O, 0, len(arr))

	for _, elem := range arr {
		var newVal = transform(elem)
		newSlice = append(newSlice, newVal)
	}

	return newSlice
}

func stringForSlice[E fmt.Stringer](elems []E, sep string) string {
	var builder strings.Builder

	var sepHasStuffInIt = len(sep) != 0
	for index, elem := range elems {
		var notTheFirstIteration = index != 0

		if notTheFirstIteration && sepHasStuffInIt {
			builder.WriteString(sep)
		}
		builder.WriteString(elem.String())
	}

	return builder.String()
}
