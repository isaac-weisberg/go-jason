package gojason

import (
	"strings"

	"github.com/isaac-weisberg/go-jason/values"
)

func debugStringAny(object *values.JsonValueAny) string {
	if object.Object != nil {
		return debugStringObject(object.Object)
	}

	if object.Number != nil {
		return debugStringNumber(object.Number)
	}

	panic("not possible (theoretically :D)")
}

func debugStringObject(object *values.JsonValueObject) string {
	var builder strings.Builder

	builder.WriteString("{")

	var kvCount = len(object.KeyValues)
	var kvIndex = 0
	for k, v := range object.KeyValues {
		builder.WriteString(debugStringAny(&k))
		builder.WriteString(":")
		builder.WriteString(debugStringAny(&v))

		var thisIsNotLastKeyValue = kvIndex != kvCount-1
		if thisIsNotLastKeyValue {
			builder.WriteString(",")
		}
		kvIndex += 1
	}
	builder.WriteString("}")

	return builder.String()
}

func debugStringNumber(number *values.JsonValueNumber) string {
	return string(number.Payload)
}
