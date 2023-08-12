package values

type JsonValueArray struct {
	values []JsonValueAny
}

func NewJsonValueArray(values []JsonValueAny) JsonValueArray {
	return JsonValueArray{
		values: values,
	}
}

func (array *JsonValueArray) AsAny() JsonValueAny {
	return JsonValueAny{
		UnderlyingType: JsonValueAnyUnderlyingTypeArray,
		Array:          array,
	}
}
