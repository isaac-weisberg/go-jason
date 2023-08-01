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

	value, err := parse(jsonString)
	if err != nil {
		t.Errorf("parse json failed with error = %v", err.Error())
	}

	fmt.Println(value.debugString())
}
