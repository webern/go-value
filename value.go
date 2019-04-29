// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Value represents a weakly typed value, like we might find in JSON. Caution: this object is not trivially copyable.
// You may use the Clone function to get a deep copy. If you use the assignment operator for copying, you will get a
// shallow copy with multiple pointers to the same address. This is almost certainly not what you want.
type Value struct {
	b    *bool
	i    *int
	f    *float64
	str  *string
	time *time.Time
	obj  Object
	arr  Array
}

func (v *Value) Equals(other Value) bool {
	t := v.Type()
	ot := other.Type()

	if t != ot {
		return false
	}

	switch t {
	case Null:
		return true
	case Bool:
		return v.Bool() == other.Bool()
	case Int:
		return v.Int() == other.Int()
	case Float:
		return v.Float() == other.Float()
	case String:
		return v.String() == other.String()
	case Time:
		return v.Time().Unix() == other.Time().Unix()
	case ArrayType:
		return ArraysEqual(v.GetArray(), other.GetArray())
	case ObjectType:
		return objectsEqual(v.GetObject(), other.GetObject())
	}

	return false
}

func (v *Value) SetType(iqType Type) {

	if iqType < Null || iqType > ArrayType {
		iqType = Null
	}

	switch iqType {
	case Null:
		{
			v.b = nil
			v.i = nil
			v.f = nil
			v.str = nil
			v.time = nil
			v.obj = nil
			v.arr = nil
		}
	case Bool:
		{
			v.b = new(bool)
			*v.b = false
			v.i = nil
			v.f = nil
			v.str = nil
			v.time = nil
			v.obj = nil
			v.arr = nil
		}
	case Int:
		{
			v.b = nil
			v.i = new(int)
			v.f = nil
			v.str = nil
			v.time = nil
			v.obj = nil
			v.arr = nil
		}
	case Float:
		{
			newD := 0.0
			v.b = nil
			v.i = nil
			v.f = &newD
			v.str = nil
			v.time = nil
			v.obj = nil
			v.arr = nil
		}
	case String:
		{
			v.b = nil
			v.i = nil
			v.f = nil
			v.str = new(string)
			v.time = nil
			v.obj = nil
			v.arr = nil
		}
	case Time:
		{
			v.b = nil
			v.i = nil
			v.f = nil
			v.str = nil
			v.time = new(time.Time)
			v.obj = nil
			v.arr = nil
		}
	case ObjectType:
		{
			v.b = nil
			v.i = nil
			v.f = nil
			v.str = nil
			v.time = nil
			v.obj = NewObject(3)
			v.arr = nil
		}
	case ArrayType:
		{
			v.b = nil
			v.i = nil
			v.f = nil
			v.str = nil
			v.time = nil
			v.obj = nil
			v.arr = NewArray()
		}
	}
}

// NewValueFromMystery makes a best effort to represent 'data' as a Value object
func NewValueFromMystery(data interface{}) (v Value, err error) {

	switch data.(type) {
	case float32:
		{
			v.SetFloat(float64(data.(float32)))
			return v, nil
		}
	case float64:
		{
			v.SetFloat(data.(float64))
			return v, nil
		}
	case int:
		{
			v.SetInt(data.(int))
			return v, nil
		}
	case int64:
		{
			v.SetInt(int(data.(int64)))
			return v, nil
		}
	case int32:
		{
			v.SetInt(int(data.(int32)))
			return v, nil
		}
	case string:
		{
			v.SetString(data.(string))
			return v, nil
		}
	case bool:
		{
			v.SetBool(data.(bool))
			return v, nil
		}
	}

	if val, ok := data.(Value); ok {
		v = val.Clone()
		return v, nil
	} else if val, ok := data.(*Value); ok {
		v = val.Clone()
		return v, nil
	} else if arr, ok := data.([]interface{}); ok {
		newArr := NewArray()

		for _, currentMystery := range arr {
			newValFromMystery, errx := NewValueFromMystery(currentMystery)

			if errx != nil {
				return v, errx
			}

			newArr = append(newArr, newValFromMystery)
		}

		v.SetArray(newArr)
		return v, nil
	} else if arr, ok := data.(*[]interface{}); ok {
		newArr := NewArray()

		for _, currentMystery := range *arr {
			newValFromMystery, errx := NewValueFromMystery(currentMystery)

			if errx != nil {
				return v, errx
			}

			newArr = append(newArr, newValFromMystery)
		}

		v.SetArray(newArr)
		return v, nil
	} else if imap, ok := data.(map[string]interface{}); ok {
		newObj := NewObject(len(imap))

		for name, currentMystery := range imap {
			newValFromMystery, errx := NewValueFromMystery(currentMystery)

			if errx != nil {
				return v, errx
			}

			newObj[name] = newValFromMystery
		}

		v.SetObject(newObj)
		return v, nil
	} else if imap, ok := data.(*map[string]interface{}); ok {
		newObj := NewObject(len(*imap))

		for name, currentMystery := range *imap {
			newValFromMystery, errx := NewValueFromMystery(currentMystery)

			if errx != nil {
				return v, errx
			}

			newObj[name] = newValFromMystery
		}

		v.SetObject(newObj)
		return v, nil
	}

	// maybe it is a struct, maybe we can marshal and unmarshal it
	b, err := json.Marshal(data)

	if err != nil {
		return v, err
	}

	err = json.Unmarshal(b, &v)

	if err != nil {
		v.SetNull()
		return v, err
	}

	return v, nil
}

func (v Value) Type() Type {

	if v.b != nil {
		return Bool
	} else if v.i != nil {
		return Int
	} else if v.f != nil {
		return Float
	} else if v.str != nil {
		return String
	} else if v.time != nil {
		return Time
	} else if v.obj != nil {
		return ObjectType
	} else if v.arr != nil {
		return ArrayType
	}

	return Null
}

func (v Value) TryBool() (value bool, err error) {
	if v.b == nil {
		return false, fmt.Errorf("TryBool was called but the type is %s", v.Type().String())
	}

	return *v.b, nil
}

func (v Value) TryInt() (value int, err error) {
	if v.i == nil {
		return 0, fmt.Errorf("TryInt was called but the type is %s", v.Type().String())
	}

	return *v.i, nil
}

func (v Value) TryFloat() (value float64, err error) {
	if v.f == nil {
		return 0.0, fmt.Errorf("TryFloat was called but the type is %s", v.Type().String())
	}

	return *v.f, nil
}

func (v Value) TryString() (value string, err error) {
	if v.str == nil {
		return "", fmt.Errorf("TryString was called but the type is %s", v.Type().String())
	}

	return *v.str, nil
}

func (v Value) TryTime() (value time.Time, err error) {
	if v.time == nil {
		return time.Time{}, fmt.Errorf("TryTime was called but the type is %s", v.Type().String())
	}

	return *v.time, nil
}

func (v Value) Time() time.Time {
	if v.time == nil {
		return time.Time{}
	}

	return *v.time
}

func (v Value) GetObject() (value Object) {
	if v.obj == nil {
		return Object{}
	}

	return v.obj
}

func (v Value) GetArray() (value Array) {
	if v.arr == nil {
		return Array{}
	}

	return v.arr
}

func (v *Value) SetNull() {
	v.SetType(Null)
}

func (v *Value) SetBool(value bool) {
	v.SetType(Bool)
	*v.b = value
}

func (v *Value) SetInt(value int) {
	v.SetType(Int)
	*v.i = value
}

func (v *Value) SetFloat(value float64) {
	v.SetType(Float)
	*v.f = value
}

func (v *Value) SetString(value string) {
	v.SetType(String)
	*v.str = value
}

func (v *Value) SetTime(value time.Time) {
	v.SetType(Time)
	*v.time = value
}

func (v *Value) SetObject(value Object) {
	v.SetType(ObjectType)
	v.obj = value
}

func (v *Value) SetArray(value Array) {
	v.SetType(ArrayType)
	v.arr = value
}

func (v Value) MarshalJSON() ([]byte, error) {

	t := v.Type()
	switch t {
	case Null:
		{
			return json.Marshal(nil)
		}
	case Bool:
		{
			base, _ := v.TryBool()
			return json.Marshal(base)
		}
	case Int:
		{
			base, _ := v.TryInt()
			return json.Marshal(base)
		}
	case Float:
		{
			base, _ := v.TryFloat()
			return json.Marshal(base)
		}
	case String:
		{
			base, _ := v.TryString()
			return json.Marshal(base)
		}
	case Time:
		{
			base, _ := v.TryTime()
			return json.Marshal(base)
		}
	case ObjectType:
		{
			base := v.GetObject()

			if base == nil {
				return make([]byte, 0), nil
			}

			return json.Marshal(base)
		}
	case ArrayType:
		{
			base := v.GetArray()

			if base == nil {
				return make([]byte, 0), nil
			}

			return json.Marshal(base)
		}
	default:
		{
			return json.Marshal(nil)
		}
	}
}

func (v *Value) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		v.SetNull()
		return nil
	}

	if len(s) >= 2 && s[0:1] == "\"" && s[len(s)-1:] == "\"" {
		str, err := strconv.Unquote(s)
		if err != nil {
			return err
		}
		v.SetString(str)
		return nil
	} else if len(s) >= 2 && s[0:1] == "{" && s[len(s)-1:] == "}" {
		o := NewObject(3)
		err := json.Unmarshal(data, &o)
		if err != nil {
			return err
		}
		v.SetObject(o)
		return nil
	} else if len(s) >= 2 && s[0:1] == "[" && s[len(s)-1:] == "]" {
		a := NewArray()
		err := json.Unmarshal(data, &a)
		if err != nil {
			return err
		}
		v.SetArray(a)
		return nil
	}

	pt := Parse(s)
	if ok := pt.Has(Null); ok {
		v.SetNull()
	} else if b, ok := pt.Bool(); ok {
		v.SetBool(b)
	} else if i, ok := pt.Integer(); ok {
		v.SetInt(i)
	} else if f, ok := pt.Float(); ok {
		v.SetFloat(f)
	} else {
		// if it's a string without any quotes?
		if len(s) > 0 {
			v.SetString(s)
		} else {
			v.SetNull()
		}
	}

	return nil
}

func (v Value) Clone() Value {
	var newVal Value

	switch v.Type() {
	case Null:
		return newVal
	case Bool:
		{
			newVal.SetBool(*v.b)
		}
	case Int:
		{
			newVal.SetInt(*v.i)
		}
	case Float:
		{
			newVal.SetFloat(*v.f)
		}
	case String:
		{
			newVal.SetString(*v.str)
		}
	case Time:
		{
			newVal.SetTime(*v.time)
		}
	case ObjectType:
		{
			newVal.SetObject(v.obj.Clone())
		}
	case ArrayType:
		{
			newVal.SetArray(v.arr.Clone())
		}
	}

	return newVal
}

func NewValue() Value {
	var val Value
	return val
}

func NewIntValue(v int) Value {
	var val Value
	val.SetInt(v)
	return val
}

func NewStringValue(v string) Value {
	var val Value
	val.SetString(v)
	return val
}

func NewBoolValue(v bool) Value {
	var val Value
	val.SetBool(v)
	return val
}

func NewFloatValue(f float64) Value {
	var val Value
	val.SetFloat(f)
	return val
}

func NewArrayValue(a Array) Value {
	var val Value
	val.SetArray(a)
	return val
}

func NewObjectValue(o Object) Value {
	v := Value{}

	if o == nil {
		v.SetObject(NewObject(3))
	} else {
		v.SetObject(o)
	}

	return v
}

func (v Value) Int() int {
	o, _ := v.TryInt()
	return o
}

func (v Value) Float() float64 {
	o, _ := v.TryFloat()
	return o
}

func (v Value) String() string {
	o, _ := v.TryString()
	return o
}

func (v Value) Bool() bool {
	o, _ := v.TryBool()
	return o
}

func (v Value) IsNull() bool {
	return v.Type() == Null
}

func (v Value) IsInt() bool {
	return v.Type() == Int
}

func (v Value) IsFloat() bool {
	return v.Type() == Float
}

func (v Value) IsString() bool {
	return v.Type() == String
}

func (v Value) IsBool() bool {
	return v.Type() == Bool
}

func (v Value) IsTime() bool {
	return v.Type() == Time
}

func (v Value) CoerceTo(t Type) (newValue Value, ok bool) {
	if t == String {
		return v.CoerceToString()
	} else if t == Int {
		return v.CoerceToInt()
	} else if t == Float {
		return v.CoearceToFloat()
	} else if t == Bool {
		return v.CoerceToBool()
	} else if t == Null {
		return Value{}, true
	}

	return Value{}, false
}

func (v Value) CoerceToInt() (newValue Value, ok bool) {
	t := v.Type()
	ok = true
	switch t {
	case Null:
		{
			newValue.SetInt(0)
			return newValue, ok
		}
	case Bool:
		{
			if newValue.Bool() {
				newValue.SetInt(1)
				return newValue, ok
			} else {
				newValue.SetInt(0)
				return newValue, ok
			}
		}
	case Int:
		return v.Clone(), ok
	case Float:
		{
			// caution, rounds to nearest int instead of truncating
			f := v.Float()
			i := int(math.Round(f))
			newValue.SetInt(i)
			return newValue, ok
		}
	case String:
		{
			s := v.String()
			i, err := strconv.Atoi(s)

			if err != nil {
				newValue.SetInt(0)
				return newValue, false
			} else {
				newValue.SetInt(i)
				return newValue, ok
			}
		}
	case Time:
		fallthrough
	case ArrayType:
		fallthrough
	case ObjectType:
		fallthrough
	default:
		break
	}

	newValue.SetInt(0)
	return newValue, false
}

func (v Value) CoearceToFloat() (newValue Value, ok bool) {
	t := v.Type()
	ok = true
	switch t {
	case Null:
		{
			newValue.SetFloat(0.0)
			return newValue, ok
		}
	case Bool:
		{
			if newValue.Bool() {
				newValue.SetFloat(1.0)
				return newValue, ok
			} else {
				newValue.SetFloat(0.0)
				return newValue, ok
			}
		}
	case Int:
		{
			i := v.Int()
			f := float64(i)
			newValue.SetFloat(f)
			return newValue, ok
		}
	case Float:
		{
			newValue = v.Clone()
			return newValue, ok
		}
	case String:
		{
			s := v.String()
			f, err := strconv.ParseFloat(s, 64)

			if err != nil {
				newValue.SetFloat(0.0)
				return newValue, false
			} else {
				newValue.SetFloat(f)
				return newValue, ok
			}
		}
	case Time:
		fallthrough
	case ArrayType:
		fallthrough
	case ObjectType:
		fallthrough
	default:
		break
	}

	newValue.SetFloat(0.0)
	return newValue, false
}

func (v Value) CoerceToString() (newValue Value, ok bool) {
	t := v.Type()
	ok = true
	if t == String {
		return v.Clone(), ok
	}

	js, err := json.Marshal(v)

	if err != nil {
		return NewStringValue(""), false
	}

	newV := NewStringValue(string(js))
	return newV, ok
}

func (v Value) CoerceToBool() (newValue Value, ok bool) {
	t := v.Type()
	ok = true
	switch t {
	case Null:
		{
			newValue.SetBool(false)
			return newValue, ok
		}
	case Bool:
		{
			return v.Clone(), ok
		}
	case Int:
		{
			i := v.Int()
			newValue.SetBool(i != 0)
			return newValue, ok
		}
	case Float:
		{
			f := v.Float()
			newValue.SetBool(math.Abs(f) > 0.0)
			return newValue, ok
		}
	case String:
		{
			s := strings.ToLower(v.String())

			if s == "true" {
				newValue.SetBool(true)
				return newValue, ok
			} else if s == "false" {
				newValue.SetBool(false)
				return newValue, ok
			} else if s == "1" {
				newValue.SetBool(true)
				return newValue, ok
			} else if s == "0" {
				newValue.SetBool(false)
				return newValue, ok
			} else if s == "yes" {
				newValue.SetBool(true)
				return newValue, ok
			} else if s == "no" {
				newValue.SetBool(false)
				return newValue, ok
			} else if s == "y" {
				newValue.SetBool(true)
				return newValue, ok
			} else if s == "n" {
				newValue.SetBool(false)
				return newValue, ok
			} else {
				newValue.SetBool(false)
				return newValue, false
			}
		}
	case Time:
		fallthrough
	case ArrayType:
		fallthrough
	case ObjectType:
		fallthrough
	default:
		break
	}

	newValue.SetBool(false)
	return newValue, false
}
