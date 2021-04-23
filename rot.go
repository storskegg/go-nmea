package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	TypeROT = "ROT"

	ValidROT   = "A"
	InvalidROT = "V"
)

// 1 	Rate of turn, degrees/minutes, “–” indicates bow turns to port
// 2 	Satus:
//			A: Valid data
// 			V: Invalid data

type ROT struct {
	BaseSentence
	RateOfTurn Float64 // Rate of turn
	Status     String
}

// newROT constructor
func newROT(s BaseSentence) (ROT, error) {
	p := NewParser(s)
	p.AssertType(TypeROT)
	m := ROT{
		BaseSentence: s,
		RateOfTurn:   p.Float64(0, "ROT"),
		Status:       p.EnumString(1, "status", ValidROT, InvalidROT),
	}
	return m, p.Err()
}

// GetRateOfTurn retrieves the rate of turn from the sentence
func (s ROT) GetRateOfTurn() (float64, error) {
	if v, err := s.RateOfTurn.GetValue(); err == nil && s.Status.Value == ValidROT {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
