package values

type JsonValueString struct {
	String string
}

func (jsonValueString *JsonValueString) AsAny() JsonValueAny {
	if jsonValueString == nil {
		panic("how")
	}
	return JsonValueAny{
		UnderlyingType: JsonValueAnyUnderlyingTypeString,
		String:         jsonValueString,
	}
}

func NewJsonValueString(value string) JsonValueString {
	return JsonValueString{
		String: value,
	}
}
