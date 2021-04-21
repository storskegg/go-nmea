package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	// TypeVHW type for VHW sentences
	TypeVHW = "VHW"
)

// VHW contains information about water speed and heading
type VHW struct {
	BaseSentence
	TrueHeading            Float64
	MagneticHeading        Float64
	SpeedThroughWaterKnots Float64
	SpeedThroughWaterKPH   Float64
}

// newVHW constructor
func newVHW(s BaseSentence) (VHW, error) {
	p := NewParser(s)
	p.AssertType(TypeVHW)
	return VHW{
		BaseSentence:           s,
		TrueHeading:            p.Float64(0, "true heading"),
		MagneticHeading:        p.Float64(2, "magnetic heading"),
		SpeedThroughWaterKnots: p.Float64(4, "speed through water in knots"),
		SpeedThroughWaterKPH:   p.Float64(6, "speed through water in kilometers per hour"),
	}, p.Err()
}

// GetMagneticHeading retrieves the magnetic heading from the sentence
func (s VHW) GetMagneticHeading() (float64, error) {
	if v, err := s.MagneticHeading.GetValue(); err == nil {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetTrueHeading retrieves the true heading from the sentence
func (s VHW) GetTrueHeading() (float64, error) {
	if v, err := s.TrueHeading.GetValue(); err == nil {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetSpeedThroughWater retrieves the speed through water from the sentence
func (s VHW) GetSpeedThroughWater() (float64, error) {
	if v, err := s.SpeedThroughWaterKPH.GetValue(); err == nil {
		return (unit.Speed(v) * unit.KilometersPerHour).MetersPerSecond(), nil
	}
	if v, err := s.SpeedThroughWaterKnots.GetValue(); err == nil {
		return (unit.Speed(v) * unit.Knot).MetersPerSecond(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
