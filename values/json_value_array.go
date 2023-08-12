package values

type JsonValueArray struct {
	Values []JsonValueAny
}

func NewJsonValueArray(values []JsonValueAny) JsonValueArray {
	return JsonValueArray{
		Values: values,
	}
}

func (array *JsonValueArray) AsAny() JsonValueAny {
	return JsonValueAny{
		UnderlyingType: JsonValueAnyUnderlyingTypeArray,
		Array:          array,
	}
}
