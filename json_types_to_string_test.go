package gojason

import "strings"

func (object *JsonValueAny) debugString() string {
	if object.object != nil {
		return object.object.debugString()
	}

	if object.number != nil {
		return object.number.debugString()
	}

	panic("not possible (theoretically :D)")
}

func (object *JsonValueObject) debugString() string {
	var builder strings.Builder

	builder.WriteString("{")

	var kvCount = len(object.keyValues)
	var kvIndex = 0
	for k, v := range object.keyValues {
		builder.WriteString(k.debugString())
		builder.WriteString(":")
		builder.WriteString(v.debugString())

		var thisIsNotLastKeyValue = kvIndex != kvCount-1
		if thisIsNotLastKeyValue {
			builder.WriteString(",")
		}
		kvIndex += 1
	}
	builder.WriteString("}")

	return builder.String()
}

func (number *JsonValueNumber) debugString() string {
	return string(number.payload)
}
