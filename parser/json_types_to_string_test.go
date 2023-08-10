package parser

import (
	"strings"

	"github.com/isaac-weisberg/go-jason/values"
)

func debugStringAny(object *values.JsonValueAny) string {
	switch object.UnderlyingType {
	case values.JsonValueAnyUnderlyingTypeNumber:
		return debugStringNumber(object.Number)
	case values.JsonValueAnyUnderlyingTypeObject:
		return debugStringObject(object.Object)
	case values.JsonValueAnyUnderlyingTypeString:
		return debugStringString(object.String)
	}

	panic("not possible (theoretically :D)")
}

func debugStringString(object *values.JsonValueString) string {
	var builder strings.Builder

	builder.WriteString("\"")
	builder.WriteString(object.String)
	builder.WriteString("\"")

	return builder.String()
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
