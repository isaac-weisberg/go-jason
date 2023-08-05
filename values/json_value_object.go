package values

type JsonValueObject struct {
	KeyValues map[JsonValueObjectKey]JsonValueAny
}

func (jsonValueObject *JsonValueObject) AsAny() JsonValueAny {
	return JsonValueAny{
		UnderlyingType: JsonValueAnyUnderlyingTypeObject,
		Object:         jsonValueObject,
	}
}

func NewJsonValueObject(keyValues map[JsonValueObjectKey]JsonValueAny) JsonValueObject {
	return JsonValueObject{
		KeyValues: keyValues,
	}
}

type JsonValueObjectKey = JsonValueAny
type JsonValueObjectKeyValues = map[JsonValueObjectKey]JsonValueAny

func NewJsonValueObjectKeyValues(cap int64) JsonValueObjectKeyValues {
	return make(JsonValueObjectKeyValues, cap)
}
