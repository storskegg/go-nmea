package nmea

import (
	"fmt"
	"time"
)

const (
	// TypeZDA type for ZDA sentences
	TypeZDA = "ZDA"
)

// ZDA represents date & time data.
// http://aprs.gids.nl/nmea/#zda
type ZDA struct {
	BaseSentence
	Time          Time
	Day           Int64
	Month         Int64
	Year          Int64
	OffsetHours   Int64 // Local time zone offset from GMT, hours
	OffsetMinutes Int64 // Local time zone offset from GMT, minutes
}

// newZDA constructor
func newZDA(s BaseSentence) (ZDA, error) {
	p := NewParser(s)
	p.AssertType(TypeZDA)
	return ZDA{
		BaseSentence:  s,
		Time:          p.Time(0, "time"),
		Day:           p.Int64(1, "day"),
		Month:         p.Int64(2, "month"),
		Year:          p.Int64(3, "year"),
		OffsetHours:   p.Int64(4, "offset (hours)"),
		OffsetMinutes: p.Int64(5, "offset (minutes)"),
	}, p.Err()
}

// GetDateTime retrieves the date and time in RFC3339Nano format
func (s ZDA) GetDateTime() (string, error) {
	if !s.Time.Valid {
		return "", fmt.Errorf("value is unavailable")
	}
	day, err := s.Day.GetValue()
	if err != nil {
		return "", fmt.Errorf("value is unavailable")
	}
	month, err := s.Month.GetValue()
	if err != nil {
		return "", fmt.Errorf("value is unavailable")
	}
	year, err := s.Year.GetValue()
	if err != nil {
		return "", fmt.Errorf("value is unavailable")
	}
	return time.Date(
		int(year),
		time.Month(month),
		int(day),
		s.Time.Hour,
		s.Time.Minute,
		s.Time.Second,
		s.Time.Millisecond*1000000,
		time.UTC,
	).UTC().Format(time.RFC3339Nano), nil
}
