package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	TypeVWR = "VWR"

	LeftOfBow  = "L"
	RightOfBow = "R"
)

// Sentence info:
// 1    Measured angle relative to the vessel, left/right of vessel heading, to the nearest 0.1 degree
// 2    L: Left, or R: Right
// 3    Measured wind speed, knots, to the nearest 0.1 knot
// 4    N: knots
// 5    Wind speed, meters per second, to the nearest 0.1 m/s
// 6    M: meters per second
// 7    Wind speed, km/h, to the nearest km/h
// 8    K: km/h

type VWR struct {
	BaseSentence
	Angle                        Float64
	LeftRightOfBow               String
	WindSpeedInKnots             Float64
	WindSpeedInMetersPerSecond   Float64
	WindSpeedInKilometersPerHour Float64
}

// newVWR constructor
func newVWR(s BaseSentence) (VWR, error) {
	p := NewParser(s)
	p.AssertType(TypeVWR)
	m := VWR{
		BaseSentence:                 s,
		Angle:                        p.Float64(0, "Angle"),
		LeftRightOfBow:               p.EnumString(1, "LeftRightOfBow", LeftOfBow, RightOfBow),
		WindSpeedInKnots:             p.Float64(2, "WindSpeedInKnots"),
		WindSpeedInMetersPerSecond:   p.Float64(4, "WindSpeedInMetersPerSecond"),
		WindSpeedInKilometersPerHour: p.Float64(6, "WindSpeedInKilometersPerHour"),
	}
	return m, p.Err()
}

// GetRelativeWindDirection retrieves the true wind direction from the sentence
func (s VWR) GetRelativeWindDirection() (float64, error) {
	if v, err := s.Angle.GetValue(); err == nil {
		if s.LeftRightOfBow.Value == LeftOfBow {
			return -(unit.Angle(v) * unit.Degree).Radians(), nil
		}
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetWindSpeed retrieves wind speed from the sentence
func (s VWR) GetWindSpeed() (float64, error) {
	if v, err := s.WindSpeedInMetersPerSecond.GetValue(); err == nil {
		return v, nil
	}
	if v, err := s.WindSpeedInKilometersPerHour.GetValue(); err == nil {
		return (unit.Speed(v) * unit.KilometersPerHour).MetersPerSecond(), nil
	}
	if v, err := s.WindSpeedInKnots.GetValue(); err == nil {
		return (unit.Speed(v) * unit.Knot).MetersPerSecond(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
