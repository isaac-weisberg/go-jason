package values

import (
	"fmt"

	"github.com/isaac-weisberg/go-jason/util"
)

type JsonValueAny struct {
	UnderlyingType JsonValueAnyUnderlyingType
	Number         *JsonValueNumber
	Object         *JsonValueObject
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=JsonValueAnyUnderlyingType -output json_types_json_value_any_underlying_type_strings.go
type JsonValueAnyUnderlyingType int64

const (
	JsonValueAnyUnderlyingTypeObject JsonValueAnyUnderlyingType = iota
	JsonValueAnyUnderlyingTypeNumber
)

func (valueAny *JsonValueAny) asObject() (*JsonValueObject, error) {
	if valueAny.Object != nil {
		return valueAny.Object, nil
	}
	return nil, util.E(fmt.Sprintf("this value can not be interpreted as an object, it has underlying type of %s", valueAny.UnderlyingType.String()))
}
