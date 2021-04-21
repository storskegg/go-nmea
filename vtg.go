package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	// TypeVTG type for VTG sentences
	TypeVTG = "VTG"
)

// VTG represents track & speed data.
// http://aprs.gids.nl/nmea/#vtg
type VTG struct {
	BaseSentence
	TrueTrack        Float64
	MagneticTrack    Float64
	GroundSpeedKnots Float64
	GroundSpeedKPH   Float64
}

// newVTG parses the VTG sentence into this struct.
// e.g: $GPVTG,360.0,T,348.7,M,000.0,N,000.0,K*43
func newVTG(s BaseSentence) (VTG, error) {
	p := NewParser(s)
	p.AssertType(TypeVTG)
	return VTG{
		BaseSentence:     s,
		TrueTrack:        p.Float64(0, "true track"),
		MagneticTrack:    p.Float64(2, "magnetic track"),
		GroundSpeedKnots: p.Float64(4, "ground speed (knots)"),
		GroundSpeedKPH:   p.Float64(6, "ground speed (km/h)"),
	}, p.Err()
}

// GetTrueCourseOverGround retrieves the true course over ground from the sentence
func (s VTG) GetTrueCourseOverGround() (float64, error) {
	if v, err := s.TrueTrack.GetValue(); err == nil {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetMagneticCourseOverGround retrieves the magnetic course over ground from the sentence
func (s VTG) GetMagneticCourseOverGround() (float64, error) {
	if v, err := s.MagneticTrack.GetValue(); err == nil {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetSpeedOverGround retrieves the speed over ground from the sentence
func (s VTG) GetSpeedOverGround() (float64, error) {
	if v, err := s.GroundSpeedKPH.GetValue(); err == nil {
		return (unit.Speed(v) * unit.KilometersPerHour).MetersPerSecond(), nil
	}
	if v, err := s.GroundSpeedKnots.GetValue(); err == nil {
		return (unit.Speed(v) * unit.Knot).MetersPerSecond(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
