package values

type JsonValueNumber struct {
	Payload []rune
}

func (jsonValueNumber *JsonValueNumber) AsAny() JsonValueAny {
	return JsonValueAny{
		UnderlyingType: JsonValueAnyUnderlyingTypeNumber,
		Number:         jsonValueNumber,
	}
}

func NewJsonValueNumber(payload []rune) JsonValueNumber {
	return JsonValueNumber{
		Payload: payload,
	}
}
