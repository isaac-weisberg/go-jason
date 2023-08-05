package gojason

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
		}
	}
	`

	value, err := Parse([]byte(jsonString))
	if err != nil {
		t.Errorf("parse json failed with error = %v", err.Error())
	}

	fmt.Println(debugStringAny(value))
}
