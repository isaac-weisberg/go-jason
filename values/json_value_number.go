package values

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
