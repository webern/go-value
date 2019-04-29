// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"encoding/json"
	"errors"
	"strings"
)

// Object represents a composite object of values, like we would find in JSON, e.g. { "someObject": { "A": 1 } }
type Object map[string]Value

// Creates a new Object which is a map[string]Value
func NewObject(capacity int) Object {
	return make(Object, capacity)
}

// Clone creates a deep copy of the object
func (o Object) Clone() Object {
	var newObj = make(Object, len(o))
	if len(o) == 0 {
		return newObj
	}

	for name, item := range o {
		newObj[name] = item.Clone()
	}
	return newObj
}

// MarshalInOrder will marshal the object to JSON placing the properties in the order that you give it.
//func (o Object) MarshalInOrder(order []string) ([]byte, error) {
//
//	b := bytes.Buffer{}
//	b.WriteString("{")
//
//	for index, key := range order {
//
//		val, ok := o[key]
//
//		if !ok {
//			continue
//		}
//
//		b.WriteString("\"")
//		b.WriteString(key)
//		b.WriteString("\":")
//		js, err := json.Marshal(val)
//		if err != nil {
//			return nil, err
//		}
//		b.Write(js)
//
//		if index < len(order)-1 {
//			b.WriteString(",")
//		}
//	}
//
//	b.WriteString("}")
//	return b.Bytes(), nil
//}

func (o *Object) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.TrimSpace(s)
	*o = NewObject(3)

	if len(s) < 2 || s[0:1] != "{" || s[len(s)-1:] != "}" {
		return errors.New("this is not a JSON object")
	}

	var temp map[string]interface{}
	err := json.Unmarshal(data, &temp)

	if err != nil {
		return err
	} else if temp == nil {
		return errors.New("bad JSON unmarshal occurred")
	}

	for name, mystery := range temp {
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
			(*o)[name] = currVal
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
			(*o)[name] = currVal
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
			(*o)[name] = currVal
		}
	}

	return nil
}

// Private

// objectsEqual returns true if all the properties and values of the two objects are the same
func objectsEqual(left, right Object) bool {
	l := len(left)
	if l != len(right) {
		return false
	}

	for k, v := range left {
		rval, ok := right[k]

		if !ok {
			return false
		}

		if !v.Equals(rval) {
			return false
		}
	}

	return true
}
