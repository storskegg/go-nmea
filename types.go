package nmea

// Latitude / longitude representation.

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	// Degrees value
	Degrees = '\u00B0'
	// Minutes value
	Minutes = '\''
	// Seconds value
	Seconds = '"'
	// Point value
	Point = '.'
	// North value
	North = "N"
	// South value
	South = "S"
	// East value
	East = "E"
	// West value
	West = "W"
)

// ParseLatLong parses the supplied string into the LatLong.
//
// Supported formats are:
// - DMS (e.g. 33° 23' 22")
// - Decimal (e.g. 33.23454)
// - GPS (e.g 15113.4322S)
//
func ParseLatLong(s string) Float64 {
	v := NewInvalidFloat64("could not parse as dms, gps or decimal notation")
	if !v.Valid {
		v = ParseDMS(s)
	}
	if !v.Valid {
		v = ParseGPS(s)
	}
	if !v.Valid {
		v = ParseDecimal(s)
	}
	if !v.Valid {
		return v
	}

	direction := string(s[len(s)-1:])
	if (direction == North || direction == South) && (v.Value < -90.0 || 90.0 < v.Value) {
		return NewInvalidFloat64("latitude is not in range (-90, 90)")
	} else if (direction == West || direction == East) && (v.Value < -180.0 || 180.0 < v.Value) {
		return NewInvalidFloat64("longitude is not in range (-180, 180)")
	}

	return v
}

// ParseGPS parses a GPS/NMEA coordinate.
// e.g 15113.4322S
func ParseGPS(s string) Float64 {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return NewInvalidFloat64(fmt.Sprintf("invalid format: %s", s))
	}
	dir := parts[1]
	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return NewInvalidFloat64(fmt.Sprintf("parse error: %s", err.Error()))
	}

	degrees := math.Floor(value / 100)
	minutes := value - (degrees * 100)
	value = degrees + minutes/60

	if dir == North || dir == East {
		return NewFloat64(value)
	}
	if dir == South || dir == West {
		return Float64{Valid: true, Value: 0 - value}
	}
	return NewInvalidFloat64(fmt.Sprintf("invalid direction [%s]", dir))
}

// FormatGPS formats a GPS/NMEA coordinate
func FormatGPS(l Float64) string {
	padding := ""
	degrees := math.Floor(math.Abs(l.Value))
	fraction := (math.Abs(l.Value) - degrees) * 60
	if fraction < 10 {
		padding = "0"
	}
	return fmt.Sprintf("%d%s%.4f", int(degrees), padding, fraction)
}

// ParseDecimal parses a decimal format coordinate.
// e.g: 151.196019
func ParseDecimal(s string) Float64 {
	// Make sure it parses as a float.
	l, err := strconv.ParseFloat(s, 64)
	if err != nil || s[0] != '-' && len(strings.Split(s, ".")[0]) > 3 {
		return NewInvalidFloat64("parse error (not decimal coordinate)")
	}
	return NewFloat64(l)
}

// ParseDMS parses a coordinate in degrees, minutes, seconds.
// - e.g. 33° 23' 22"
func ParseDMS(s string) Float64 {
	degrees := 0
	minutes := 0
	seconds := 0.0
	// Whether a number has finished parsing (i.e whitespace after it)
	endNumber := false
	// Temporary parse buffer.
	tmpBytes := []byte{}
	var err error

	for i, r := range s {
		switch {
		case unicode.IsNumber(r) || r == '.':
			if !endNumber {
				tmpBytes = append(tmpBytes, s[i])
			} else {
				return NewInvalidFloat64("parse error (no delimiter)")
			}
		case unicode.IsSpace(r) && len(tmpBytes) > 0:
			endNumber = true
		case r == Degrees:
			if degrees, err = strconv.Atoi(string(tmpBytes)); err != nil {
				return NewInvalidFloat64("parse error (degrees)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		case s[i] == Minutes:
			if minutes, err = strconv.Atoi(string(tmpBytes)); err != nil {
				return NewInvalidFloat64("parse error (minutes)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		case s[i] == Seconds:
			if seconds, err = strconv.ParseFloat(string(tmpBytes), 64); err != nil {
				return NewInvalidFloat64("parse error (seconds)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		case unicode.IsSpace(r) && len(tmpBytes) == 0:
			continue
		default:
			return NewInvalidFloat64(fmt.Sprintf("parse error (unknown symbol [%d])", s[i]))
		}
	}
	if len(tmpBytes) > 0 {
		return NewInvalidFloat64(fmt.Sprintf("parse error (trailing data [%s])", string(tmpBytes)))
	}
	val := float64(degrees) + (float64(minutes) / 60.0) + (float64(seconds) / 60.0 / 60.0)
	return NewFloat64(val)
}

// FormatDMS returns the degrees, minutes, seconds format for the given LatLong.
func FormatDMS(l Float64) string {
	val := math.Abs(l.Value)
	degrees := int(math.Floor(val))
	minutes := int(math.Floor(60 * (val - float64(degrees))))
	seconds := 3600 * (val - float64(degrees) - (float64(minutes) / 60))
	return fmt.Sprintf("%d\u00B0 %d' %f\"", degrees, minutes, seconds)
}

// Time type
type Time struct {
	Valid         bool
	InvalidReason string
	Hour          int
	Minute        int
	Second        int
	Millisecond   int
}

func NewTime(hour int, minute int, second int, millisecond int) Time {
	return Time{
		Valid:       true,
		Hour:        hour,
		Minute:      minute,
		Second:      second,
		Millisecond: millisecond,
	}
}

func NewInvalidTime(reason string) Time {
	return Time{
		Valid:         false,
		InvalidReason: reason,
	}
}

// String representation of Time
func (t Time) String() string {
	seconds := float64(t.Second) + float64(t.Millisecond)/1000
	return fmt.Sprintf("%02d:%02d:%07.4f", t.Hour, t.Minute, seconds)
}

// timeRe is used to validate time strings
var timeRe = regexp.MustCompile(`^\d{6}(\.\d*)?$`)

// ParseTime parses wall clock time.
// e.g. hhmmss.ssss
// An empty time string will result in an invalid time.
func ParseTime(s string) Time {
	if !timeRe.MatchString(s) {
		return NewInvalidTime(fmt.Sprintf("parse time: expected hhmmss.ss format, got '%s'", s))
	}
	hour, _ := strconv.Atoi(s[:2])
	minute, _ := strconv.Atoi(s[2:4])
	second, _ := strconv.ParseFloat(s[4:], 64)
	whole, frac := math.Modf(second)
	return NewTime(hour, minute, int(whole), int(math.Round(frac*1000)))
}

// Date type
type Date struct {
	Valid         bool
	InvalidReason string
	DD            int
	MM            int
	YY            int
}

func NewDate(year int, month int, day int) Date {
	return Date{
		Valid: true,
		YY:    year,
		MM:    month,
		DD:    day,
	}
}

func NewInvalidDate(reason string) Date {
	return Date{
		Valid:         false,
		InvalidReason: reason,
	}
}

// String representation of date
func (d Date) String() string {
	return fmt.Sprintf("%02d/%02d/%02d", d.DD, d.MM, d.YY)
}

// dateRe is used to validate date strings
var dateRe = regexp.MustCompile(`^\d{6}$`)

// ParseDate field ddmmyy format
func ParseDate(s string) Date {
	if !dateRe.MatchString(s) {
		return NewInvalidDate(fmt.Sprintf("parse date: expected yymmdd format, got '%s'", s))
	}
	dd, _ := strconv.Atoi(s[0:2])
	mm, _ := strconv.Atoi(s[2:4])
	yy, _ := strconv.Atoi(s[4:6])
	return NewDate(yy, mm, dd)
}

// LatDir returns the latitude direction symbol
func LatDir(l float64) string {
	if l < 0.0 {
		return South
	}
	return North
}

// LonDir returns the longitude direction symbol
func LonDir(l float64) string {
	if l < 0.0 {
		return East
	}
	return West
}

type Float64 struct {
	Valid         bool
	InvalidReason string
	Value         float64
}

func NewFloat64(v float64) Float64 {
	return Float64{
		Valid: true,
		Value: v,
	}
}

func NewInvalidFloat64(reason string) Float64 {
	return Float64{
		InvalidReason: reason,
	}
}

func (v Float64) GetValue() (float64, error) {
	if v.Valid {
		return v.Value, nil
	}
	return 0, errors.New(v.InvalidReason)
}

func ParseFloat64(s string) Float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return NewInvalidFloat64(err.Error())
	}
	return NewFloat64(v)
}

type Int64 struct {
	Valid         bool
	InvalidReason string
	Value         int64
}

func NewInt64(v int64) Int64 {
	return Int64{
		Valid: true,
		Value: v,
	}
}
func NewInvalidInt64(reason string) Int64 {
	return Int64{
		InvalidReason: reason,
	}
}

func (v Int64) GetValue() (int64, error) {
	if v.Valid {
		return v.Value, nil
	}
	return 0, errors.New(v.InvalidReason)
}

func ParseInt64(s string) Int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return NewInvalidInt64(err.Error())
	}
	return NewInt64(v)
}

type String struct {
	Valid         bool
	InvalidReason string
	Value         string
}

func NewString(v string) String {
	return String{
		Valid: true,
		Value: v,
	}
}

func NewInvalidString(reason string) String {
	return String{
		InvalidReason: reason,
	}
}

func (v String) GetValue() (string, error) {
	if v.Valid {
		return v.Value, nil
	}
	return "", fmt.Errorf("the value is invalid")
}

type StringList struct {
	Valid         bool
	InvalidReason string
	Values        []String
}

func NewStringList(v []String) StringList {
	return StringList{
		Valid:  true,
		Values: v,
	}
}

func NewInvalidStringList(reason string) StringList {
	return StringList{
		InvalidReason: reason,
	}
}
