// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"bytes"
	"encoding/json"
)

// Type represents the value types that we can have in
type Type int

const (
	Null       Type = iota // Null has bi value
	Bool                   // Bool holds a boolean value
	Int                    // Int holds an int value
	Float                  // Float holds a float64 value
	String                 // String holds a string value
	Time                   // Time holds a time.Time value
	ObjectType             // ObjectType holds an Object, which is a map[string]Value
	ArrayType              // ArrayType holds a slice which is []Value
)

const (
	StringNull    = "VALUE_NULL"
	StringBool    = "VALUE_BOOL"
	StringInt     = "VALUE_INTEGER"
	StringDecimal = "VALUE_DECIMAL"
	StringString  = "VALUE_STRING"
	StringTime    = "VALUE_TIME"
	StringObject  = "VALUE_OBJECT"
	StringArray   = "VALUE_ARRAY"
)

var typeToString = map[Type]string{
	Null:       StringNull,
	Bool:       StringBool,
	Int:        StringInt,
	Float:      StringDecimal,
	String:     StringString,
	Time:       StringTime,
	ObjectType: StringObject,
	ArrayType:  StringArray,
}

var stringToType = map[string]Type{
	StringNull:    Null,
	StringBool:    Bool,
	StringInt:     Int,
	StringDecimal: Float,
	StringString:  String,
	StringTime:    Time,
	StringObject:  ObjectType,
	StringArray:   ArrayType,
}

func (t Type) String() string {

	theString, ok := typeToString[t]

	if ok {
		return theString
	}

	return typeToString[Null]
}

func (t *Type) Parse(s string) {

	theEnumValue, ok := stringToType[s]

	if ok {
		*t = theEnumValue
	} else {
		*t = Null
	}
}

func (t Type) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	_, err := buffer.WriteString(t.String())

	if err != nil {
		return buffer.Bytes(), err
	}

	_, err = buffer.WriteString("\"")

	if err != nil {
		return buffer.Bytes(), err
	}

	return buffer.Bytes(), nil
}

func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	t.Parse(s)
	return nil
}
