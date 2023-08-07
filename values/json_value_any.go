package values

import (
	"github.com/isaac-weisberg/go-jason/util"
)

type JsonValueAny struct {
	UnderlyingType JsonValueAnyUnderlyingType
	Number         *JsonValueNumber
	Object         *JsonValueObject
	String         *JsonValueString
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=JsonValueAnyUnderlyingType -output json_types_json_value_any_underlying_type_strings.go
type JsonValueAnyUnderlyingType int64

const (
	JsonValueAnyUnderlyingTypeObject JsonValueAnyUnderlyingType = iota
	JsonValueAnyUnderlyingTypeNumber
	JsonValueAnyUnderlyingTypeString
)

const valueInterpretationErrorBaseString = "this value can not be interpreted as an object, it has underlying type of %s"

func (valueAny *JsonValueAny) AsObject() (*JsonValueObject, error) {
	if valueAny.Object != nil {
		return valueAny.Object, nil
	}
	return nil, util.E(valueInterpretationErrorBaseString, valueAny.UnderlyingType.String())
}

func (valueAny *JsonValueAny) AsString() (*JsonValueString, error) {
	if valueAny.String != nil {
		return valueAny.String, nil
	}
	return nil, util.E(valueInterpretationErrorBaseString, valueAny.UnderlyingType.String())
}

func (valueAny *JsonValueAny) AsNumber() (*JsonValueNumber, error) {
	if valueAny.Number != nil {
		return valueAny.Number, nil
	}
	return nil, util.E(valueInterpretationErrorBaseString, valueAny.UnderlyingType.String())
}
