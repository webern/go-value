// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import "testing"

func TestNewValueFromMystery(t *testing.T) {
	type TestCase struct {
		Mystery       interface{}
		Expected      Value
		IsErrExpected bool
	}

	testCases := []TestCase{
		{
			Mystery: struct {
				Hello string
				World int
			}{Hello: "Blerp", World: 5},
			Expected: NewObjectValue(Object{
				"Hello": NewStringValue("Blerp"),
				"World": NewIntValue(5),
			}),
			IsErrExpected: false,
		},
	}

	for tcix, tc := range testCases {

	}
}
