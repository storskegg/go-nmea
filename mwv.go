package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	TypeMWV = "MWV"

	ReferenceRelative = "R"
	ReferenceTrue     = "T"

	WindSpeedUnitKPH   = "K"
	WindSpeedUnitKnots = "N"
	WindSpeedUnitMPS   = "M"
	WindSpeedUnitMPH   = "S"

	ValidMWV   = "A"
	InvalidMWV = "V"
)

// Sentence info:
// 1    Wind angle, 0.0 to 359.9 degrees, in relation to the vessel’s bow/centerline, to the nearest 0.1 degree. If the data for this field is not valid, the field will be blank.
// 2    Reference:
// 			R: Relative (apparent wind, as felt when standing on the moving ship)
// 			T: Theoretical (calculated actual wind, as though the vessel were stationary)
// 3    Wind speed, to the nearest tenth of a unit.  If the data for this field is not valid, the field will be blank.
// 4    Wind speed units:
//			K: km/hr
//			M: m/s
//			N: knots
//			S: statute miles/hr
// 5    Status:
//			A: data valid
//			V: data invalid

type MWV struct {
	BaseSentence
	Angle         Float64
	Reference     string
	WindSpeed     Float64
	WindSpeedUnit string
	Status        string
}

// newMWV constructor
func newMWV(s BaseSentence) (MWV, error) {
	p := NewParser(s)
	p.AssertType(TypeMWV)
	m := MWV{
		BaseSentence:  s,
		Angle:         p.Float64(0, "Angle"),
		Reference:     p.EnumString(1, "Reference", ReferenceRelative, ReferenceTrue),
		WindSpeed:     p.Float64(2, "WindSpeed"),
		WindSpeedUnit: p.EnumString(3, "WindSpeedUnit", WindSpeedUnitKPH, WindSpeedUnitKnots, WindSpeedUnitMPS, WindSpeedUnitMPH),
		Status:        p.EnumString(4, "Status", ValidMWV, InvalidMWV),
	}
	return m, p.Err()
}

// GetTrueWindDirection retrieves the true wind direction from the sentence
func (s MWV) GetTrueWindDirection() (float64, error) {
	if s.Status == ValidMWV && s.Reference == ReferenceTrue {
		if v, err := s.Angle.GetValue(); err == nil {
			return (unit.Angle(v) * unit.Degree).Radians(), nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetRelativeWindDirection retrieves the relative wind direction from the sentence
func (s MWV) GetRelativeWindDirection() (float64, error) {
	if s.Status == ValidMWV && s.Reference == ReferenceRelative {
		if v, err := s.Angle.GetValue(); err == nil {
			return (unit.Angle(v) * unit.Degree).Radians(), nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetWindSpeed retrieves wind speed from the sentence
func (s MWV) GetWindSpeed() (float64, error) {
	if v, err := s.WindSpeed.GetValue(); err == nil && s.Status == ValidMWV {
		switch s.WindSpeedUnit {
		case WindSpeedUnitMPS:
			return v, nil
		case WindSpeedUnitKPH:
			return (unit.Speed(v) * unit.KilometersPerHour).MetersPerSecond(), nil
		case WindSpeedUnitKnots:
			return (unit.Speed(v) * unit.Knot).MetersPerSecond(), nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}
