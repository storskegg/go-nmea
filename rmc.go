package nmea

import (
	"fmt"
	"time"

	"github.com/martinlindhe/unit"
)

const (
	// TypeRMC type for RMC sentences
	TypeRMC = "RMC"
	// ValidRMC character
	ValidRMC = "A"
	// InvalidRMC character
	InvalidRMC = "V"
)

// RMC is the Recommended Minimum Specific GNSS data.
// http://aprs.gids.nl/nmea/#rmc
type RMC struct {
	BaseSentence
	Time      Time    // Time Stamp
	Validity  String  // validity - A-ok, V-invalid
	Latitude  Float64 // Latitude
	Longitude Float64 // Longitude
	Speed     Float64 // Speed in knots
	Course    Float64 // True course
	Date      Date    // Date
	Variation Float64 // Magnetic variation
}

// newRMC constructor
func newRMC(s BaseSentence) (RMC, error) {
	p := NewParser(s)
	p.AssertType(TypeRMC)
	m := RMC{
		BaseSentence: s,
		Time:         p.Time(0, "time"),
		Validity:     p.EnumString(1, "validity", ValidRMC, InvalidRMC),
		Latitude:     p.LatLong(2, 3, "latitude"),
		Longitude:    p.LatLong(4, 5, "longitude"),
		Speed:        p.Float64(6, "speed"),
		Course:       p.Float64(7, "course"),
		Date:         p.Date(8, "date"),
		Variation:    p.Float64(9, "variation"),
	}
	if m.Variation.Valid && p.EnumString(10, "direction", West, East).Value == West {
		m.Variation.Value = 0 - m.Variation.Value
	}
	return m, p.Err()
}

// GetMagneticVariation retrieves the magnetic variation from the sentence
func (s RMC) GetMagneticVariation() (float64, error) {
	if s.Validity.Value == ValidRMC {
		if v, err := s.Variation.GetValue(); err == nil {
			return (unit.Angle(v) * unit.Degree).Radians(), nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetTrueCourseOverGround retrieves the true course over ground from the sentence
func (s RMC) GetTrueCourseOverGround() (float64, error) {
	if s.Validity.Value == ValidRMC {
		if v, err := s.Course.GetValue(); err == nil {
			return (unit.Angle(v) * unit.Degree).Radians(), nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetPosition2D retrieves the latitude and longitude from the sentence
func (s RMC) GetPosition2D() (float64, float64, error) {
	if s.Validity.Value == ValidRMC {
		if latitude, err := s.Latitude.GetValue(); err == nil {
			if longitude, err := s.Longitude.GetValue(); err == nil {
				return latitude, longitude, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("value is unavailable")
}

// GetSpeedOverGround retrieves the speed over ground from the sentence
func (s RMC) GetSpeedOverGround() (float64, error) {
	if s.Validity.Value == ValidRMC {
		if v, err := s.Speed.GetValue(); err == nil {
			return (unit.Speed(v) * unit.Knot).MetersPerSecond(), nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetDateTime retrieves the date and time in RFC3339Nano format
func (s RMC) GetDateTime() (string, error) {
	if s.Validity.Value == ValidRMC {
		if s.Date.Valid && s.Time.Valid {
			return time.Date(
				s.Date.YY,
				time.Month(s.Date.MM),
				s.Date.DD,
				s.Time.Hour,
				s.Time.Minute,
				s.Time.Second,
				s.Time.Millisecond*1000000,
				time.UTC,
			).UTC().Format(time.RFC3339Nano), nil
		}
	}
	return "", fmt.Errorf("value is unavailable")
}
