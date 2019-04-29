// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"fmt"
	"testing"
	"time"

	"github.com/webern/tcore"
)

func TestNewValueFromMystery(t *testing.T) {
	type TestCase struct {
		Mystery       interface{}
		Expected      Value
		IsErrExpected bool
	}

	someTime, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon May 6 10:00:00 -0700 MST 2019")
	someVal := NewStringValue("flerbin")

	if err != nil {
		t.Error(err.Error())
		return
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
		{
			Mystery:       5,
			Expected:      NewIntValue(5),
			IsErrExpected: false,
		},
		{
			Mystery:       "flerbin",
			Expected:      NewStringValue("flerbin"),
			IsErrExpected: false,
		},
		{
			Mystery:       1.01,
			Expected:      NewFloatValue(1.01),
			IsErrExpected: false,
		},
		{
			Mystery:       someTime,
			Expected:      NewTimeValue(someTime),
			IsErrExpected: false,
		},
		{
			Mystery:       someTime,
			Expected:      NewTimeValue(someTime),
			IsErrExpected: false,
		},
		{
			Mystery:       float32(2.0),
			Expected:      NewFloatValue(2.0),
			IsErrExpected: false,
		},
		{
			Mystery:       int32(3),
			Expected:      NewIntValue(3),
			IsErrExpected: false,
		},
		{
			Mystery:       int64(4),
			Expected:      NewIntValue(4),
			IsErrExpected: false,
		},
		{
			Mystery:       true,
			Expected:      NewBoolValue(true),
			IsErrExpected: false,
		},
		{
			Mystery:       false,
			Expected:      NewBoolValue(false),
			IsErrExpected: false,
		},
		{
			Mystery:       someVal,
			Expected:      someVal.Clone(),
			IsErrExpected: false,
		},
		{
			Mystery:       &someVal,
			Expected:      someVal.Clone(),
			IsErrExpected: false,
		},
		{
			Mystery:       []int{1, 2},
			Expected:      NewArrayValue(Array{NewIntValue(1), NewIntValue(2)}),
			IsErrExpected: false,
		},
	}

	for tcix, tc := range testCases {
		stm := fmt.Sprintf("test case %d: %s", tcix, "actual, err := NewValueFromMystery(tc.Mystery)")
		actual, err := NewValueFromMystery(tc.Mystery)

		if tc.IsErrExpected {
			if err == nil {
				t.Errorf("%s - an error was expected but none was raised", stm)
			}
		} else {
			if msg, ok := tcore.TErr(stm, err); !ok {
				t.Error(msg)
			}
		}

		stm = fmt.Sprintf("test case %d: %s", tcix, "actual.Equals(tc.Expected)")
		gotB := actual.Equals(tc.Expected)
		if msg, ok := tcore.TAssertBool(stm, gotB, true); !ok {
			t.Error(msg)
		}
	}
}
