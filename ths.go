package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	// TypeTHS type for THS sentences
	TypeTHS = "THS"
	// AutonomousTHS autonomous ths heading
	AutonomousTHS = "A"
	// EstimatedTHS estimated (dead reckoning) THS heading
	EstimatedTHS = "E"
	// ManualTHS manual input THS heading
	ManualTHS = "M"
	// SimulatorTHS simulated THS heading
	SimulatorTHS = "S"
	// InvalidTHS not valid THS heading (or standby)
	InvalidTHS = "V"
)

// THS is the Actual vessel heading in degrees True with status.
// http://www.nuovamarea.net/pytheas_9.html
type THS struct {
	BaseSentence
	Heading Float64 // Heading in degrees
	Status  String  // Heading status
}

// newTHS constructor
func newTHS(s BaseSentence) (THS, error) {
	p := NewParser(s)
	p.AssertType(TypeTHS)
	m := THS{
		BaseSentence: s,
		Heading:      p.Float64(0, "heading"),
		Status:       p.EnumString(1, "status", AutonomousTHS, EstimatedTHS, ManualTHS, SimulatorTHS, InvalidTHS),
	}
	return m, p.Err()
}

// GetTrueHeading retrieves the true heading from the sentence
func (s THS) GetTrueHeading() (float64, error) {
	if v, err := s.Heading.GetValue(); err == nil && s.Status.Value != InvalidTHS {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
