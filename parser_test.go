package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	var (
		sentence string
		p        *Parser
	)
	Describe("Testing methods of parser", func() {
		JustBeforeEach(func() {
			parsed, _ := Parse(sentence)
			if typedParsed, ok := parsed.(RMC); ok {
				p = NewParser(typedParsed.BaseSentence)
				return
			}
			if typedParsed, ok := parsed.(GNS); ok {
				p = NewParser(typedParsed.BaseSentence)
				return
			}
			if typedParsed, ok := parsed.(RTE); ok {
				p = NewParser(typedParsed.BaseSentence)
				return
			}
			Fail("This type of sentence is not supported")
		})
		Context("when setting an error", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns the error", func() {
				p.SetErr("testing", "error")
				Expect(p.Err()).To(MatchError("nmea: GNRMC invalid testing: error"))
			})
		})
		Context("when asserting it is a RMC sentence", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns no error", func() {
				p.AssertType("RMC")
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when asserting it is a ROT sentence", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an error", func() {
				p.AssertType("ROT")
				Expect(p.Err()).To(HaveOccurred())
			})
		})
		Context("when retrieving a string field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns the value of the field", func() {
				Expect(p.String(6, "speed")).To(Equal(NewString("0.036")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a string field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid string", func() {
				Expect(p.String(23, "invalid index")).To(Equal(NewInvalidString("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a list of string field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$IIRTE,X,1,c,Rte 1,411,412,413,414,415*03"
			})
			It("returns the value of the field", func() {
				Expect(p.ListString(4, "ident of waypoints")).To(Equal(NewStringList([]String{NewString("411"), NewString("412"), NewString("413"), NewString("414"), NewString("415")})))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a list of string field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$IIRTE,X,1,c,Rte 1,411,412,413,414,415*03"
			})
			It("returns an invalid string", func() {
				Expect(p.ListString(23, "invalid index")).To(Equal(NewInvalidStringList("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an enum string field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns the value of the field", func() {
				Expect(p.EnumString(10, "direction", West, East)).To(Equal(NewString(East)))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an enum string field with an invalid value", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,X,A*2C"
			})
			It("returns the value of the field", func() {
				Expect(p.EnumString(10, "direction", West, East)).To(Equal(NewInvalidString("not a valid option")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an enum string field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid string", func() {
				Expect(p.EnumString(23, "invalid index", West, East)).To(Equal(NewInvalidString("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an enum char field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,AA,14,0.6,161.5,48.0,,*6D"
			})
			It("returns the value of the field", func() {
				Expect(p.EnumChars(5, "mode", NoFixGNS, AutonomousGNS, DifferentialGNS, PreciseGNS, RealTimeKinematicGNS, FloatRTKGNS, EstimatedGNS, ManualGNS, SimulatorGNS)).To(Equal(NewStringList([]String{NewString("A"), NewString("A")})))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an enum char field with an invalid value", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,AAX,14,0.6,161.5,48.0,,*35"
			})
			It("returns the value of the field", func() {
				Expect(p.EnumChars(5, "mode", NoFixGNS, AutonomousGNS, DifferentialGNS, PreciseGNS, RealTimeKinematicGNS, FloatRTKGNS, EstimatedGNS, ManualGNS, SimulatorGNS)).To(Equal(NewStringList([]String{NewString("A"), NewString("A"), NewInvalidString("not a valid option")})))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an enum char field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid string", func() {
				Expect(p.EnumChars(23, "invalid index", West, East)).To(Equal(NewInvalidStringList("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an int64 field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,PXKR,14,0.6,161.5,48.0,,*7C"
			})
			It("returns the value of the field", func() {
				Expect(p.Int64(6, "satellites")).To(Equal(NewInt64(14)))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an int64 field with an invalid value", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,PXKR,X,0.6,161.5,48.0,,*21"
			})
			It("returns the value of the field", func() {
				Expect(p.Int64(6, "satellites")).To(Equal(NewInvalidInt64("strconv.ParseInt: parsing \"X\": invalid syntax")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an int64 field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,PXKR,14,0.6,161.5,48.0,,*7C"
			})
			It("returns an invalid int64", func() {
				Expect(p.Int64(23, "invalid index")).To(Equal(NewInvalidInt64("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an float64 field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,PXKR,14,0.6,161.5,48.0,,*7C"
			})
			It("returns the value of the field", func() {
				Expect(p.Float64(8, "altitude")).To(Equal(NewFloat64(161.5)))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an float64 field with an invalid value", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,PXKR,14,0.6,X,48.0,,*09"
			})
			It("returns the value of the field", func() {
				Expect(p.Float64(8, "altitude")).To(Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"X\": invalid syntax")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving an float64 field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,PXKR,14,0.6,161.5,48.0,,*7C"
			})
			It("returns an invalid float64", func() {
				Expect(p.Float64(23, "invalid index")).To(Equal(NewInvalidFloat64("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a time field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns the value of the field", func() {
				Expect(p.Time(0, "time")).To(Equal(NewTime(14, 39, 9, 0)))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a time field with an invalid value", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,X,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*41"
			})
			It("returns the value of the field", func() {
				Expect(p.Time(0, "time")).To(Equal(NewInvalidTime("parse time: expected hhmmss.ss format, got 'X'")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a time field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid time", func() {
				Expect(p.Time(23, "invalid index")).To(Equal(NewInvalidTime("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a date field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns the value of the field", func() {
				Expect(p.Date(8, "date")).To(Equal(NewDate(7, 3, 21)))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a date field with an invalid value", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,X,A,5107.0020216,N,11402.3294835,W,0.036,348.3,X,0.0,E,A*1E"
			})
			It("returns the value of the field", func() {
				Expect(p.Date(8, "date")).To(Equal(NewInvalidDate("parse date: expected yymmdd format, got 'X'")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a date field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid date", func() {
				Expect(p.Date(23, "invalid index")).To(Equal(NewInvalidDate("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a latitude field with a valid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns the value of the field", func() {
				Expect(p.LatLong(2, 3, "latitude")).To(Equal(NewFloat64(51.11670036)))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a position field with an invalid value", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,X,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*73"
			})
			It("returns the value of the field", func() {
				Expect(p.LatLong(2, 3, "latitude")).To(Equal(NewInvalidFloat64("parse error (not decimal coordinate)")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a position field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid date", func() {
				Expect(p.LatLong(2, 24, "invalid index")).To(Equal(NewInvalidFloat64("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a position field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid date", func() {
				Expect(p.LatLong(23, 3, "invalid index")).To(Equal(NewInvalidFloat64("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
		Context("when retrieving a position field with an invalid index", func() {
			BeforeEach(func() {
				sentence = "$GNRMC,143909.00,A,5107.0020216,N,11402.3294835,W,0.036,348.3,210307,0.0,E,A*31"
			})
			It("returns an invalid date", func() {
				Expect(p.LatLong(23, 24, "invalid index")).To(Equal(NewInvalidFloat64("index out of range")))
				Expect(p.Err()).ToNot(HaveOccurred())
			})
		})
	})
})

// import (
// 	"fmt"
// 	"testing"

// 	. "github.com/munnik/go-nmea"
// 	"github.com/stretchr/testify/assert"
// )

// var parsertests = []struct {
// 	name     string
// 	fields   []string
// 	expected interface{}
// 	hasErr   bool
// 	parse    func(p *Parser) interface{}
// }{
// 	{
// 		name:   "Bad Type",
// 		fields: []string{},
// 		hasErr: true,
// 		parse: func(p *Parser) interface{} {
// 			p.AssertType("WRONG_TYPE")
// 			return nil
// 		},
// 	},
// 	{
// 		name:     "String",
// 		fields:   []string{"foo", "bar"},
// 		expected: "bar",
// 		parse: func(p *Parser) interface{} {
// 			return p.String(1, "")
// 		},
// 	},
// 	{
// 		name:     "String out of range",
// 		fields:   []string{"wot"},
// 		expected: "",
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.String(5, "thing")
// 		},
// 	},
// 	{
// 		name:     "ListString",
// 		fields:   []string{"wot", "foo", "bar"},
// 		expected: []string{"foo", "bar"},
// 		parse: func(p *Parser) interface{} {
// 			return p.ListString(1, "thing")
// 		},
// 	},
// 	{
// 		name:     "ListString out of range",
// 		fields:   []string{"wot"},
// 		expected: []string{},
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.ListString(10, "thing")
// 		},
// 	},
// 	{
// 		name:     "String with existing error",
// 		expected: "",
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.String(123, "blah")
// 		},
// 	},
// 	{
// 		name:     "EnumString",
// 		fields:   []string{"a", "b", "c"},
// 		expected: "b",
// 		parse: func(p *Parser) interface{} {
// 			return p.EnumString(1, "context", "b", "d")
// 		},
// 	},
// 	{
// 		name:     "EnumString invalid",
// 		fields:   []string{"a", "b", "c"},
// 		expected: "",
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.EnumString(1, "context", "x", "y")
// 		},
// 	},
// 	{
// 		name:     "EnumString with existing error",
// 		fields:   []string{"a", "b", "c"},
// 		expected: "",
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.EnumString(1, "context", "a", "b")
// 		},
// 	},
// 	{
// 		name:     "EnumChars",
// 		fields:   []string{"AA", "AB", "BA", "BB"},
// 		expected: []string{"A", "B"},
// 		parse: func(p *Parser) interface{} {
// 			return p.EnumChars(1, "context", "A", "B")
// 		},
// 	},
// 	{
// 		name:     "EnumChars invalid",
// 		fields:   []string{"a", "AB", "c"},
// 		expected: []string{},
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.EnumChars(1, "context", "X", "Y")
// 		},
// 	},
// 	{
// 		name:     "EnumChars with existing error",
// 		fields:   []string{"a", "AB", "c"},
// 		expected: []string{},
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.EnumChars(1, "context", "A", "B")
// 		},
// 	},
// 	{
// 		name:     "Int64",
// 		fields:   []string{"123"},
// 		expected: NewInt64(123),
// 		parse: func(p *Parser) interface{} {
// 			return p.Int64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Int64 empty field is zero",
// 		fields:   []string{""},
// 		expected: NewInvalidInt64(errors.New("")),
// 		parse: func(p *Parser) interface{} {
// 			return p.Int64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Int64 invalid",
// 		fields:   []string{"abc"},
// 		expected: NewInvalidInt64(errors.New("")),
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.Int64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Int64 with existing error",
// 		fields:   []string{"123"},
// 		expected: NewInvalidInt64(errors.New("")),
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.Int64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Float64",
// 		fields:   []string{"123.123"},
// 		expected: NewFloat64(123.123),
// 		parse: func(p *Parser) interface{} {
// 			return p.Float64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Float64 empty field is zero",
// 		fields:   []string{""},
// 		expected: NewInvalidFloat64(errors.New("")),
// 		parse: func(p *Parser) interface{} {
// 			return p.Float64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Float64 invalid",
// 		fields:   []string{"abc"},
// 		expected: NewInvalidFloat64(errors.New("")),
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.Float64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Float64 with existing error",
// 		fields:   []string{"123.123"},
// 		expected: NewInvalidFloat64(errors.New("")),
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.Float64(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Time",
// 		fields:   []string{"123456"},
// 		expected: NewTime(12, 34, 56, 0),
// 		parse: func(p *Parser) interface{} {
// 			return p.Time(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Time empty field is zero",
// 		fields:   []string{""},
// 		expected: Time{},
// 		parse: func(p *Parser) interface{} {
// 			return p.Time(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Time with existing error",
// 		fields:   []string{"123456"},
// 		expected: Time{},
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.Time(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Time invalid",
// 		fields:   []string{"wrong"},
// 		expected: Time{},
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.Time(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Date",
// 		fields:   []string{"010203"},
// 		expected: NewDate(3, 2, 1),
// 		parse: func(p *Parser) interface{} {
// 			return p.Date(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Date empty field is zero",
// 		fields:   []string{""},
// 		expected: Date{},
// 		parse: func(p *Parser) interface{} {
// 			return p.Date(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Date invalid",
// 		fields:   []string{"Hello"},
// 		expected: Date{},
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			return p.Date(0, "context")
// 		},
// 	},
// 	{
// 		name:     "Date with existing error",
// 		fields:   []string{"010203"},
// 		expected: Date{},
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.Date(0, "context")
// 		},
// 	},
// 	{
// 		name:     "LatLong",
// 		fields:   []string{"5000.0000", "N"},
// 		expected: NewFloat64(50.0),
// 		parse: func(p *Parser) interface{} {
// 			return p.LatLong(0, 1, "context")
// 		},
// 	},
// 	{
// 		name:     "LatLong - latitude out of range",
// 		fields:   []string{"9100.0000", "N"},
// 		expected: NewInvalidFloat64(fmt.Errorf("latitude is not in range (-90, 90)")),
// 		hasErr:   false,
// 		parse: func(p *Parser) interface{} {
// 			return p.LatLong(0, 1, "context")
// 		},
// 	},
// 	{
// 		name:     "LatLong - longitude out of range",
// 		fields:   []string{"18100.0000", "W"},
// 		expected: NewInvalidFloat64(fmt.Errorf("longitude is not in range (-180, 180)")),
// 		hasErr:   false,
// 		parse: func(p *Parser) interface{} {
// 			return p.LatLong(0, 1, "context")
// 		},
// 	},
// 	{
// 		name:     "LatLong with existing error",
// 		fields:   []string{"5000.0000", "W"},
// 		expected: NewInvalidFloat64(errors.New("")),
// 		hasErr:   true,
// 		parse: func(p *Parser) interface{} {
// 			p.SetErr("context", "value")
// 			return p.LatLong(0, 1, "context")
// 		},
// 	},
// }

// func TestParser(t *testing.T) {
// 	for _, tt := range parsertests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewParser(BaseSentence{
// 				Talker: "talker",
// 				Type:   "type",
// 				Fields: tt.fields,
// 			})
// 			assert.Equal(t, tt.expected, tt.parse(p))
// 			if tt.hasErr {
// 				assert.Error(t, p.Err())
// 			} else {
// 				assert.NoError(t, p.Err())
// 			}
// 		})
// 	}
// }
