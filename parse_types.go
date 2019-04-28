// go-value, Copyright (c) 2019 by Matthew James Briggs

package value

//type ParseType int

//const (
//	Nothing    ParseType = iota // it is an empty string
//	Null                        // the string may represent a null value, e.g. 'null' or 'NULL'
//	Bool                        // the string may represent 'true' or 'false'
//	Integer                     // the string may represent an integer, e.g. no decimal point is present
//	Float                     // the string may represent a decimal, e.g. a decimal point is present
//	Scientific                  // the string may represent scientific notation, e.g. a decimal point and 'E' or 'e' with exponent are present
//)
//
//const (
//	StringNothing    = "STR_NOTHING"
//	StringNull       = "STR_NULL"
//	StringInteger    = "STR_INTEGER"
//	StringDecimal    = "STR_DECIMAL"
//	StringScientific = "STR_SCIENTIFIC"
//)
//
//var parseTypeToString = map[ParseType]string{
//	Nothing:    StringNothing,
//	Null:       StringNull,
//	Integer:    StringInteger,
//	Float:    StringDecimal,
//	Scientific: StringScientific,
//}
//
//var stringToParseType = map[string]ParseType{
//	StringNothing:    Nothing,
//	StringNull:       Null,
//	StringInteger:    Integer,
//	StringDecimal:    Float,
//	StringScientific: Scientific,
//}
//
//func (p ParseType) String() string {
//
//	theString, ok := parseTypeToString[p]
//
//	if ok {
//		return theString
//	}
//
//	return parseTypeToString[Nothing]
//}
//
//func (p *ParseType) Parse(s string) {
//
//	theEnumValue, ok := stringToParseType[s]
//
//	if ok {
//		*p = theEnumValue
//	} else {
//		*p = Nothing
//	}
//}
//
//func (p ParseType) MarshalJSON() ([]byte, error) {
//	buffer := bytes.NewBufferString("\"")
//	buffer.WriteString(p.String())
//	buffer.WriteString("\"")
//	return buffer.Bytes(), nil
//}
//
//func (p *ParseType) UnmarshalJSON(b []byte) error {
//	var s string
//	err := json.Unmarshal(b, &s)
//	if err != nil {
//		return err
//	}
//	p.Parse(s)
//	return nil
//}
