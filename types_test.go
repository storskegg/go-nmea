package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Describe("Testing methods of types", func() {
		Context("when parsing latitude or longitude", func() {
			It("returns a valid latitude or longitude", func() {
				result, err := ParseLatLong("33\u00B0 12' 34.3423\"").GetValue()
				Expect(result).To(Float64Equal(33.209540, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns a valid latitude or longitude", func() {
				result, err := ParseLatLong("3345.1232 N").GetValue()
				Expect(result).To(Float64Equal(33.752054, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns a valid latitude or longitude", func() {
				result, err := ParseLatLong("151.234532").GetValue()
				Expect(result).To(Float64Equal(151.234532, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseLatLong("151.234.532").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("parse error (not decimal coordinate)"))
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseLatLong("14035.1232 N").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("latitude is not in range (-90, 90)"))
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseLatLong("-24035.1232 E").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("longitude is not in range (-180, 180)"))
			})
		})
		Context("when parsing gps", func() {
			It("returns a valid latitude or longitude", func() {
				result, err := ParseGPS("3345.1232 N").GetValue()
				Expect(result).To(Float64Equal(33.752054, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns a valid latitude or longitude", func() {
				result, err := ParseGPS("15145.9877 S").GetValue()
				Expect(result).To(Float64Equal(-151.76646, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseGPS("12345.1234 X").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("invalid direction [X]"))
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseGPS("1234.1234").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("invalid format: 1234.1234"))
			})
		})
		Context("when parsing dms", func() {
			It("returns a valid latitude or longitude", func() {
				result, err := ParseDMS("33\u00B0 12' 34.3423\"").GetValue()
				Expect(result).To(Float64Equal(33.209540, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseDMS("33\u00B0 1.1' 34.3423\"").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("parse error (minutes)"))
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseDMS("3.3\u00B0 1' 34.3423\"").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("parse error (degrees)"))
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseDMS("33\u00B0 1' 34.34.23\"").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("parse error (seconds)"))
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseDMS("33\u00B0 1' 34.34.23\"").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("parse error (seconds)"))
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseDMS("123").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("parse error (trailing data [123])"))
			})
		})
		Context("when parsing decimal", func() {
			It("returns a valid latitude or longitude", func() {
				result, err := ParseDecimal("151.234532").GetValue()
				Expect(result).To(Float64Equal(151.234532, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns a valid latitude or longitude", func() {
				result, err := ParseDecimal("-151.234532").GetValue()
				Expect(result).To(Float64Equal(-151.234532, 0.00001))
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns an invalid latitude or longitude", func() {
				result, err := ParseDecimal("-151.234532 N").GetValue()
				Expect(result).To(BeZero())
				Expect(err).To(MatchError("parse error (not decimal coordinate)"))
			})
		})
		Context("when parsing time", func() {
			It("returns a valid time", func() {
				result := ParseTime("123456")
				Expect(result).To(Equal(NewTime(12, 34, 56, 0)))
			})
			It("returns an invalid time", func() {
				result := ParseTime("")
				Expect(result).To(Equal(NewInvalidTime("parse time: expected hhmmss.ss format, got ''")))
			})
			It("returns a valid time", func() {
				result := ParseTime("112233.123")
				Expect(result).To(Equal(NewTime(11, 22, 33, 123)))
			})
			It("returns a valid time", func() {
				result := ParseTime("010203.04")
				Expect(result).To(Equal(NewTime(1, 2, 3, 40)))
			})
			It("returns an invalid time", func() {
				result := ParseTime("10203.04")
				Expect(result).To(Equal(NewInvalidTime("parse time: expected hhmmss.ss format, got '10203.04'")))
			})
			It("returns an invalid time", func() {
				result := ParseTime("xx2233.123")
				Expect(result).To(Equal(NewInvalidTime("parse time: expected hhmmss.ss format, got 'xx2233.123'")))
			})
			It("returns an invalid time", func() {
				result := ParseTime("11xx33.123")
				Expect(result).To(Equal(NewInvalidTime("parse time: expected hhmmss.ss format, got '11xx33.123'")))
			})
			It("returns an invalid time", func() {
				result := ParseTime("1122xx.123")
				Expect(result).To(Equal(NewInvalidTime("parse time: expected hhmmss.ss format, got '1122xx.123'")))
			})
			It("returns an invalid time", func() {
				result := ParseTime("112233.xxx")
				Expect(result).To(Equal(NewInvalidTime("parse time: expected hhmmss.ss format, got '112233.xxx'")))
			})
		})
		Context("when parsing date", func() {
			It("returns a valid date", func() {
				result := ParseDate("010203")
				Expect(result).To(Equal(NewDate(3, 2, 1)))
			})
			It("returns an invalid date", func() {
				result := ParseDate("01003")
				Expect(result).To(Equal(NewInvalidDate("parse date: expected yymmdd format, got '01003'")))
			})
			It("returns an invalid date", func() {
				result := ParseDate("")
				Expect(result).To(Equal(NewInvalidDate("parse date: expected yymmdd format, got ''")))
			})
			It("returns an invalid date", func() {
				result := ParseDate("xx0203")
				Expect(result).To(Equal(NewInvalidDate("parse date: expected yymmdd format, got 'xx0203'")))
			})
			It("returns an invalid date", func() {
				result := ParseDate("01xx03")
				Expect(result).To(Equal(NewInvalidDate("parse date: expected yymmdd format, got '01xx03'")))
			})
			It("returns an invalid date", func() {
				result := ParseDate("0102xx")
				Expect(result).To(Equal(NewInvalidDate("parse date: expected yymmdd format, got '0102xx'")))
			})
		})
		Context("when formatting a latitude or longitude", func() {
			It("returns a valid string", func() {
				value := NewFloat64(151.434367)
				Expect(FormatDMS(value)).To(Equal("151° 26' 3.721200\""))
				Expect(FormatGPS(value)).To(Equal("15126.0620"))
			})
			It("returns a valid string", func() {
				value := NewFloat64(33.94057166666666)
				Expect(FormatDMS(value)).To(Equal("33° 56' 26.058000\""))
				Expect(FormatGPS(value)).To(Equal("3356.4343"))
			})
			It("returns a valid string", func() {
				value := NewFloat64(45.0)
				Expect(FormatDMS(value)).To(Equal("45° 0' 0.000000\""))
				Expect(FormatGPS(value)).To(Equal("4500.0000"))
			})
		})
		Context("when formatting a time", func() {
			It("returns a valid string", func() {
				value := NewTime(1, 2, 3, 4)
				Expect(value.String()).To(Equal("01:02:03.0040"))
			})
		})
		Context("when formatting a date", func() {
			It("returns a valid string", func() {
				value := NewDate(3, 2, 1)
				Expect(value.String()).To(Equal("01/02/03"))
			})
		})
		Context("when getting the direction of a latitude", func() {
			It("returns the correct direction", func() {
				Expect(LatDir(50.0)).To(Equal(North))
				Expect(LatDir(-50.0)).To(Equal(South))
			})
		})
		Context("when getting the direction of a longitude", func() {
			It("returns the correct direction", func() {
				Expect(LonDir(100.0)).To(Equal(West))
				Expect(LonDir(-100.0)).To(Equal(East))
			})
		})
	})
})
