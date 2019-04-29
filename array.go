// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"encoding/json"
	"errors"
	"strings"
)

// Array represents a JSON array (a slice of Values)
type Array []Value

// Creates a new Array, which is a slice of Values
func NewArray() Array {
	return make(Array, 0, 10)
}

// Clone creates a deep copy of the Array
func (a Array) Clone() Array {
	var newArr = make(Array, len(a))

	if len(a) == 0 {
		return newArr
	}

	for x, item := range a {
		newArr[x] = item.Clone()
	}

	return newArr
}

func (a *Array) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.TrimSpace(s)
	*a = NewArray()

	if len(s) < 2 || s[0:1] != "[" || s[len(s)-1:] != "]" {
		return errors.New("this is not a JSON array")
	}

	var temp []interface{}
	err := json.Unmarshal(data, &temp)

	if err != nil {
		return err
	} else if temp == nil {
		return errors.New("bad JSON unmarshal occurred")
	}

	for _, mystery := range temp {
		if imap, ok := mystery.(map[string]interface{}); ok {
			// the sub-value is another object
			currVal := Value{}
			subObj := NewObject(len(imap))
			currBytes, err := json.Marshal(imap)
			if err != nil {
				return err
			}
			err = json.Unmarshal(currBytes, &subObj)
			if err != nil {
				return err
			}
			currVal.SetObject(subObj)
			*a = append(*a, currVal)
		} else if iarr, ok := mystery.([]interface{}); ok {
			// the sub-value is an array
			currVal := Value{}
			subArr := NewArray()
			currBytes, err := json.Marshal(iarr)
			if err != nil {
				return err
			}
			err = json.Unmarshal(currBytes, &subArr)
			if err != nil {
				return err
			}
			currVal.SetArray(subArr)
			*a = append(*a, currVal)
		} else {
			currVal := Value{}
			currBytes, err := json.Marshal(mystery)
			if err != nil {
				return err
			}
			err = json.Unmarshal(currBytes, &currVal)
			if err != nil {
				return err
			}
			*a = append(*a, currVal)
		}
	}

	return nil
}

// Private

// ArraysEqual returns true if the arrays are of the same length and all values are equal (in the same order)
func ArraysEqual(left, right Array) bool {
	l := len(left)
	if l != len(right) {
		return false
	}

	for i := 0; i < l; i++ {
		if !left[i].Equals(right[i]) {
			return false
		}
	}

	return true
}
