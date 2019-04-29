// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"encoding/json"
	"fmt"
	"testing"
)
import "github.com/webern/tcore"

func TestNewArray(t *testing.T) {
	a := NewArray()
	stm := "a != nil"
	gotB := a != nil

	if msg, ok := tcore.TAssertBool(stm, gotB, true); !ok {
		t.Error(msg)
	}

	zeroA := NewArray()
	zeroB := zeroA.Clone()

	if &zeroA == &zeroB || len(zeroB) != 0 {
		t.Error("you lose, the clone function did not do what we wanted it to")
	}
}

func TestArray_Clone(t *testing.T) {
	a := NewArray()
	a = append(a, NewIntValue(1))
	a = append(a, NewStringValue("hi"))
	b := a.Clone()

	stm := "&a == &b"
	gotB := &a == &b

	if msg, ok := tcore.TAssertBool(stm, gotB, false); !ok {
		t.Error(msg)
	}

	stm = "len(b)"
	gotI := len(b)

	if msg, ok := tcore.TAssertInt(stm, gotI, 2); !ok {
		t.Error(msg)
	}

	stm = "b[0].Int()"
	gotI = b[0].Int()

	if msg, ok := tcore.TAssertInt(stm, gotI, 1); !ok {
		t.Error(msg)
	}

	stm = "b[1].String()"
	gotS := b[1].String()

	if msg, ok := tcore.TAssertString(stm, gotS, "hi"); !ok {
		t.Error(msg)
	}
}

func TestArray_UnmarshalJSON(t *testing.T) {
	type TestCase struct {
		Input           string
		IsErrorExpected bool
		ExpectedArray   Array
	}

	testCases := []TestCase{
		{
			Input:           "[]",
			IsErrorExpected: false,
			ExpectedArray:   NewArray(),
		},
		{
			Input:           `["Hello","World"]`,
			IsErrorExpected: false,
			ExpectedArray:   Array{NewStringValue("Hello"), NewStringValue("World")},
		},
		{
			Input:           `[1,true,2.1]`,
			IsErrorExpected: false,
			ExpectedArray:   Array{NewIntValue(1), NewBoolValue(true), NewFloatValue(2.1)},
		},
		{
			Input:           `[1.000000000000000000001,1.99999999999999]`,
			IsErrorExpected: false,
			ExpectedArray:   Array{NewIntValue(1), NewIntValue(2)},
		},
		{
			Input:           `[null,"null",false]`,
			IsErrorExpected: false,
			ExpectedArray:   Array{NewValue(), NewStringValue("null"), NewBoolValue(false)},
		},
		{
			Input:           `]`,
			IsErrorExpected: true,
			ExpectedArray:   NewArray(),
		},
		{
			Input:           `[xxxx]`,
			IsErrorExpected: true,
			ExpectedArray:   NewArray(),
		},
		{
			Input:           `[["x"]]`,
			IsErrorExpected: false,
			ExpectedArray:   Array{NewArrayValue(Array{NewStringValue("x")})},
		},
		{
			Input:           `[{ "hello": "world" }]`,
			IsErrorExpected: false,
			ExpectedArray: Array{
				NewObjectValue(Object{
					"hello": NewStringValue("world"),
				}),
			},
		},
	}

	for tcix, tc := range testCases {
		actual := NewArray()
		stm := fmt.Sprintf("test case %d: err := json.Unmarshal([]byte(tc.Input), &actual)", tcix)
		err := json.Unmarshal([]byte(tc.Input), &actual)

		if tc.IsErrorExpected {
			if err == nil {
				t.Errorf("an error was expected but none was received for the statement '%s'", stm)
			}

			// make sure our own function catches the error too
			stm = fmt.Sprintf("test case %d: actual.UnmarshalJSON([]byte(tc.Input))", tcix)
			err = actual.UnmarshalJSON([]byte(tc.Input))
			if err == nil {
				t.Errorf("an error was expected but none was received for the statement '%s'", stm)
			}

		} else {
			if msg, ok := tcore.TErr(stm, err); !ok {
				t.Error(msg)
			}
		}

		stm = fmt.Sprintf("test case %d: ArraysEqual(actual, tc.ExpectedArray)", tcix)
		gotB := ArraysEqual(actual, tc.ExpectedArray)
		if msg, ok := tcore.TAssertBool(stm, gotB, true); !ok {
			t.Error(msg + fmt.Sprintf("' - %v' != '%v'", actual, tc.ExpectedArray))
		}
	}
}
