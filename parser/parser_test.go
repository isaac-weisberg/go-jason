package parser

import (
	"fmt"
	"testing"
)

func TestParseSimpleJson(t *testing.T) {
	jsonString := `
{
	{
		35: 60
	}: {
		90: 35
	},
	"something": {
		"keyarooni": {
			355: "fuck"
		},
		"steven": [
			"beeven",
			"heeven",
			3,
			{
				"i said": 3
			}
		]
	}
}
`

	value, err := Parse([]byte(jsonString))
	if err != nil {
		t.Errorf("parse json failed with error = %v", err.Error())
	}

	if value != nil {
		fmt.Println(debugStringAny(value))
	}
}
