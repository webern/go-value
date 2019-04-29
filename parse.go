// go-value, Copyright (c) 2019-present by Matthew James Briggs

package value

import (
	"math"
	"strconv"
	"strings"
)

const epsilon = 0.00000000000001
const StringFalse = "false"
const StringTrue = "true"

// ParseResult is a list of the possible ways that the string can be parsed, along with the resultant parsed value
//type ParseResult map[ParseType]interface{}
type ParseResult struct {
	IsNull   bool
	IsBool   bool
	IsInt    bool
	IsFloat  bool
	TheBool  bool
	TheInt   int
	TheFloat float64
}

// Has returns true if the map contains the given ParseType
func (p ParseResult) Has(pt Type) bool {
	switch pt {
	//case Nothing:
	//	return p.IsNothing
	case Null:
		return p.IsNull
	case Bool:
		return p.IsBool
	case Int:
		return p.IsInt
	case Float:
		return p.IsFloat
	default:
		break
	}

	return false
}

func (p ParseResult) Integer() (value int, ok bool) {
	if p.IsInt {
		return p.TheInt, true
	}
	return 0, false
}

func (p ParseResult) Float() (value float64, ok bool) {
	if p.IsFloat {
		return p.TheFloat, true
	}
	return 0, false
}

func (p ParseResult) Bool() (value bool, ok bool) {
	if p.IsBool {
		return p.TheBool, true
	}
	return false, false
}

func Parse(s string) ParseResult {
	pt := ParseResult{}
	if isNullString(&s) {
		pt.IsNull = true
	} else if isTrueString(&s) {
		pt.IsBool = true
		pt.TheBool = true
	} else if isFalseString(&s) {
		pt.IsBool = true
		pt.TheBool = false
	} else if isNumericPeek(&s) {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			pt.IsFloat = true
			pt.TheFloat = f
			i := int(f)
			f2 := float64(i)
			diff := math.Abs(f - f2)
			if diff < epsilon {
				pt.IsInt = true
				pt.TheInt = i
			} else {
				i++
				f2 = float64(i)
				diff = math.Abs(f - f2)
				if diff < epsilon {
					pt.IsInt = true
					pt.TheInt = int(i)
				}
			}
		}
	}

	return pt
}

func isNullString(s *string) bool {
	if len(*s) < 3 || len(*s) > 4 {
		return false
	}

	if *s == "null" || *s == "nil" {
		return true
	}

	l := strings.ToLower(*s)

	if l == "null" || l == "nil" {
		return true
	}

	return false
}

func isTrueString(s *string) bool {
	if len(*s) != 4 {
		return false
	}

	if *s == StringTrue {
		return true
	}

	l := strings.ToLower(*s)

	ret := l == StringTrue
	return ret
}

func isFalseString(s *string) bool {
	if len(*s) != 5 {
		return false
	}

	if *s == StringFalse {
		return true
	}

	l := strings.ToLower(*s)

	ret := l == StringFalse
	return ret
}

// false means the string is not a number. true, means that the string *might* be a number
func isNumericPeek(s *string) bool {
	if len(*s) == 0 {
		return false
	}

	var first rune
	for _, r := range *s {
		first = r
		break
	}

	if first != '-' &&
		first != '.' &&
		first != '0' &&
		first != '1' &&
		first != '2' &&
		first != '3' &&
		first != '4' &&
		first != '5' &&
		first != '6' &&
		first != '7' &&
		first != '8' &&
		first != '9' {
		return false
	}

	return true
}
