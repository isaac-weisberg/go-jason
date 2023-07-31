package gojason

import "strconv"

type NumberJsonValue struct {
	intValue   *int64
	floatValue *float64
}

func ParseNumberJsonValue(payload []rune) (*NumberJsonValue, error) {
	var stringValue = string(payload)
	var intValue int64

	intValue, intErr := strconv.ParseInt(stringValue, 10, 64)
	if intErr == nil {
		return &NumberJsonValue{
			intValue:   &intValue,
			floatValue: nil,
		}, nil
	}

	var floatValue float64

	floatValue, floatErr := strconv.ParseFloat(stringValue, 64)

	if floatErr == nil {
		return &NumberJsonValue{
			intValue:   nil,
			floatValue: &floatValue,
		}, nil
	}

	return nil, j(e("parsing int and float failed"), floatErr, intErr)
}
