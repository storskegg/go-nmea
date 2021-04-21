package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	// TypeDBT type for DBT sentences
	TypeDBT = "DBT"
)

// DBT - Depth below transducer
// https://gpsd.gitlab.io/gpsd/NMEA.html#_dbt_depth_below_transducer
type DBT struct {
	BaseSentence
	DepthFeet    Float64
	DepthMeters  Float64
	DepthFathoms Float64
}

// newDBT constructor
func newDBT(s BaseSentence) (DBT, error) {
	p := NewParser(s)
	p.AssertType(TypeDBT)
	return DBT{
		BaseSentence: s,
		DepthFeet:    p.Float64(0, "depth_feet"),
		DepthMeters:  p.Float64(2, "depth_meters"),
		DepthFathoms: p.Float64(4, "depth_fathoms"),
	}, p.Err()
}

// GetDepthBelowTransducer retrieves the depth below the transducer from the sentence
func (s DBT) GetDepthBelowTransducer() (float64, error) {
	if v, err := s.DepthMeters.GetValue(); err == nil {
		return v, nil
	}
	if v, err := s.DepthFeet.GetValue(); err == nil {
		return (unit.Length(v) * unit.Foot).Meters(), nil
	}
	if v, err := s.DepthFathoms.GetValue(); err == nil {
		return (unit.Length(v) * unit.Fathom).Meters(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
