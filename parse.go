// go-value, Copyright (c) 2019 by Matthew James Briggs

package value

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const epsilon = 0.00000000000001

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
				pt.TheInt = int(i)
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

	if *s == "true" {
		return true
	}

	l := strings.ToLower(*s)

	if l == "true" {
		return true
	}

	return false
}

func isFalseString(s *string) bool {
	if len(*s) != 5 {
		return false
	}

	if *s == "false" {
		return true
	}

	l := strings.ToLower(*s)

	if l == "false" {
		return true
	}

	return false
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

// you should first use isNumericPeek to make sure s begins with a valid char
func possibleNumerics(s *string) (isInt, isDecimal, isScientific bool) {
	isInt = true
	for ix, r := range *s {
		if ix == 0 && r == '-' {
			continue
		} else if r == '.' {
			isInt = false
			isDecimal = true
		} else if r == 'e' {
			isInt = false
			isScientific = true
			isDecimal = false
		} else if r == 'E' {
			isInt = false
			isScientific = true
			isDecimal = false
		} else if r < '0' {
			return false, false, false
		} else if r > '9' {
			return false, false, false
		}
	}

	return isInt, isDecimal, isScientific
}

func ParseScientific(s *string) (value float64, ok bool) {
	value = 0.0
	isNegative := false
	isExponentEncountered := false
	isBeginningOfExponent := false
	isExponentNegative := false
	isBeginning := true
	left := bytes.Buffer{}
	right := bytes.Buffer{}

	for _, r := range *s {
		if !isExponentEncountered {
			if isBeginning && r == '-' {
				isNegative = true
				isBeginning = false
			} else if unicode.IsNumber(r) || r == '.' {
				isBeginning = false
				left.WriteRune(r)
			} else if r == 'e' || r == 'E' {
				isExponentEncountered = true
				isBeginningOfExponent = true
			} else {
				return 0.0, false
			}
		} else if isExponentEncountered {
			if isBeginningOfExponent && r == '-' {
				isExponentNegative = true
				isBeginningOfExponent = false
			} else if unicode.IsNumber(r) {
				isBeginningOfExponent = false
				right.WriteRune(r)
			} else {
				return 0.0, false
			}
		}
	}

	leftStr := left.String()
	rightStr := right.String()
	leftFloat, err := strconv.ParseFloat(leftStr, 64)

	if err != nil {
		return 0.0, false
	}

	rightInt64, err := strconv.ParseInt(rightStr, 10, strconv.IntSize)
	rightInt := int(rightInt64)

	if isNegative {
		leftFloat *= -1
	}

	if isExponentNegative {
		rightInt *= -1
	}

	pow := math.Pow10(rightInt)
	value = leftFloat * pow
	return value, true
}
