package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	// TypeRSA type for RSA sentences
	TypeRSA = "RSA"

	ValidRSA   = "A"
	InvalidRSA = "V"
)

// 1 	Rudder Angle (single or starboard rudder), degrees, “–” indicates rudder is turned to port
// 2 	Status:
//			A: Valid data
// 			V: Invalid data
// 3 	Rudder Angle (portside rudder), degrees, “–” indicates rudder is turned to port
// 4 	Status:
//			A: Valid data
// 			V: Invalid data

// RSA - Rudder Angle
type RSA struct {
	BaseSentence
	RudderAngleStarboard Float64
	StatusStarboard      String
	RudderAnglePortside  Float64
	StatusPortside       String
}

// newRSA constructor
func newRSA(s BaseSentence) (RSA, error) {
	p := NewParser(s)
	p.AssertType(TypeRSA)
	m := RSA{
		BaseSentence:         s,
		RudderAngleStarboard: p.Float64(0, "rudder angle starboard"),
		StatusStarboard:      p.EnumString(1, "status starboard", ValidRSA, InvalidRSA),
		RudderAnglePortside:  p.Float64(2, "rudder angle portside"),
		StatusPortside:       p.EnumString(3, "status portside", ValidRSA, InvalidRSA),
	}
	return m, p.Err()
}

// GetRudder retrieves the rudder angle of the single rudder from the sentence
func (s RSA) GetRudderAngle() (float64, error) {
	if _, err := s.GetRudderAnglePortside(); err == nil {
		return 0, fmt.Errorf("not a single rudder system, use the specific functions for the startboard and portside rudder")
	}
	return s.GetRudderAngleStarboard()
}

// GetRudderAngleStarboard retrieves the rudder angle of the startboard rudder from the sentence
func (s RSA) GetRudderAngleStarboard() (float64, error) {
	if v, err := s.RudderAngleStarboard.GetValue(); err == nil && s.StatusStarboard.Value == ValidRSA {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetRudderAnglePortside retrieves the rudder angle of the portside rudder from the sentence
func (s RSA) GetRudderAnglePortside() (float64, error) {
	if v, err := s.RudderAnglePortside.GetValue(); err == nil && s.StatusPortside.Value == ValidRSA {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
