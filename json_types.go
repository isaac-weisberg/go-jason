package gojason

type JsonValueNumber struct {
	payload []rune
}

func (jsonValueNumber *JsonValueNumber) asAny() JsonValueAny {
	return JsonValueAny{
		number: jsonValueNumber,
	}
}

func newJsonValueNumber(payload []rune) JsonValueNumber {
	return JsonValueNumber{
		payload: payload,
	}
}

type JsonValueObject struct {
	values map[JsonValueObjectKey]JsonValueAny
}

func (jsonValueObject *JsonValueObject) asAny() JsonValueAny {
	return JsonValueAny{
		object: jsonValueObject,
	}
}

func newJsonValueObject(keyValues map[JsonValueObjectKey]JsonValueAny) JsonValueObject {
	return JsonValueObject{
		values: keyValues,
	}
}

type JsonValueAny struct {
	number *JsonValueNumber
	object *JsonValueObject
}

type JsonValueObjectKey = JsonValueAny
type JsonValueObjectKeyValues = map[JsonValueObjectKey]JsonValueAny

func newJsonValueObjectKeyValues() JsonValueObjectKeyValues {
	return make(JsonValueObjectKeyValues)
}
