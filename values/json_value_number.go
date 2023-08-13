package values

import (
	"fmt"
	"strconv"

	"github.com/isaac-weisberg/go-jason/util"
)

type JsonValueNumber struct {
	Payload []byte
}

func (jsonValueNumber *JsonValueNumber) AsAny() JsonValueAny {
	return JsonValueAny{
		UnderlyingType: JsonValueAnyUnderlyingTypeNumber,
		Number:         jsonValueNumber,
	}
}

func NewJsonValueNumber(payload []byte) JsonValueNumber {
	return JsonValueNumber{
		Payload: payload,
	}
}

func (jsonValueNumber *JsonValueNumber) ParseInt64() (int64, error) {
	var stringPayload = string(jsonValueNumber.Payload)
	var i, err = strconv.ParseInt(stringPayload, 10, 64)
	if err != nil {
		return 0, util.W(err, fmt.Sprintf("parse int on '%s' failed", stringPayload))
	}

	return i, nil
}
