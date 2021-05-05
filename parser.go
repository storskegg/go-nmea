package nmea

import (
	"fmt"
)

// Parser provides a simple way of accessing and parsing
// sentence fields
type Parser struct {
	BaseSentence
	err error
}

// NewParser constructor
func NewParser(s BaseSentence) *Parser {
	return &Parser{BaseSentence: s}
}

// AssertType makes sure the sentence's type matches the provided one.
func (p *Parser) AssertType(typ string) {
	if p.Type != typ {
		p.SetErr("type", p.Type)
	}
}

// Err returns the first error encountered during the parser's usage.
func (p *Parser) Err() error {
	return p.err
}

// SetErr assigns an error. Calling this method has no
// effect if there is already an error.
func (p *Parser) SetErr(context, value string) {
	if p.err == nil {
		p.err = fmt.Errorf("nmea: %s invalid %s: %s", p.Prefix(), context, value)
	}
}

// String returns the field value at the specified index.
func (p *Parser) String(i int, context string) String {
	if i < 0 || i >= len(p.Fields) {
		return NewInvalidString("index out of range")
	}
	return NewString(p.Fields[i])
}

// ListString returns a list of all fields from the given start index.
// An error occurs if there is no fields after the given start index.
func (p *Parser) ListString(from int, context string) StringList {
	if from < 0 || from >= len(p.Fields) {
		return NewInvalidStringList("index out of range")
	}
	result := make([]String, 0)
	for i := from; i < len(p.Fields); i++ {
		result = append(result, p.String(i, context))
	}
	return NewStringList(result)
}

// EnumString returns the field value at the specified index.
// An error occurs if the value is not one of the options and not empty.
func (p *Parser) EnumString(i int, context string, options ...string) String {
	s := p.String(i, context)
	if !s.Valid {
		return s
	}
	for _, o := range options {
		if o == s.Value {
			return s
		}
	}
	return NewInvalidString("not a valid option")
}

// EnumChars returns an array of strings that are matched in the Mode field.
// It will only match the number of characters that are in the Mode field.
// If the value is empty, it will return an empty array
func (p *Parser) EnumChars(i int, context string, options ...string) StringList {
	s := p.String(i, context)
	if !s.Valid {
		return NewInvalidStringList(s.InvalidReason)
	}
	result := make([]String, 0)
	for _, r := range s.Value {
		rs := string(r)
		found := false
		for _, o := range options {
			if o == rs {
				result = append(result, NewString(o))
				found = true
				break
			}
		}
		if !found {
			result = append(result, NewInvalidString("not a valid option"))
		}
	}
	return NewStringList(result)
}

// Int64 returns the int64 value at the specified index.
// If the value is an empty string, 0 is returned.
func (p *Parser) Int64(i int, context string) Int64 {
	s := p.String(i, context)
	if !s.Valid {
		return NewInvalidInt64(s.InvalidReason)
	}
	return ParseInt64(s.Value)
}

// Float64 returns the float64 value at the specified index.
// If the value is an empty string, 0 is returned.
func (p *Parser) Float64(i int, context string) Float64 {
	s := p.String(i, context)
	if !s.Valid {
		return NewInvalidFloat64(s.InvalidReason)
	}
	return ParseFloat64(s.Value)
}

// Time returns the Time value at the specified index.
// If the value is empty, the Time is marked as invalid.
func (p *Parser) Time(i int, context string) Time {
	s := p.String(i, context)
	if !s.Valid {
		return NewInvalidTime(s.InvalidReason)
	}
	return ParseTime(s.Value)
}

// Date returns the Date value at the specified index.
// If the value is empty, the Date is marked as invalid.
func (p *Parser) Date(i int, context string) Date {
	s := p.String(i, context)
	if !s.Valid {
		return NewInvalidDate(s.InvalidReason)
	}
	return ParseDate(s.Value)
}

// LatLong returns the coordinate value of the specified fields.
func (p *Parser) LatLong(i, j int, context string) Float64 {
	a := p.String(i, context)
	if !a.Valid {
		return NewInvalidFloat64(a.InvalidReason)
	}
	b := p.String(j, context)
	if !b.Valid {
		return NewInvalidFloat64(b.InvalidReason)
	}
	s := fmt.Sprintf("%s %s", a.Value, b.Value)
	v := ParseLatLong(s)

	return v
}

// SixBitASCIIArmour decodes the 6-bit ascii armor used for VDM and VDO messages
func (p *Parser) SixBitASCIIArmour(i int, fillBits int, context string) []byte {
	if p.err != nil {
		return nil
	}
	if fillBits < 0 || fillBits >= 6 {
		p.SetErr(context, "fill bits")
		return nil
	}

	payload := []byte(p.String(i, "encoded payload").Value)
	numBits := len(payload)*6 - fillBits

	if numBits < 0 {
		p.SetErr(context, "num bits")
		return nil
	}

	result := make([]byte, numBits)
	resultIndex := 0

	for _, v := range payload {
		if v < 48 || v >= 120 {
			p.SetErr(context, "data byte")
			return nil
		}

		d := v - 48
		if d > 40 {
			d -= 8
		}

		for i := 5; i >= 0 && resultIndex < len(result); i-- {
			result[resultIndex] = (d >> uint(i)) & 1
			resultIndex++
		}
	}

	return result
}
