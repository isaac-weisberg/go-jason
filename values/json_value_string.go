package values

type JsonValueString struct {
	Payload []byte
}

func (jsonValueString *JsonValueString) AsAny() JsonValueAny {
	return JsonValueAny{
		UnderlyingType: JsonValueAnyUnderlyingTypeString,
		String:         jsonValueString,
	}
}

func NewJsonValueString(payload []byte) JsonValueString {
	return JsonValueString{
		Payload: payload,
	}
}
