package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

const (
	// TypeDBS type for DBS sentences
	TypeDBS = "DBS"
)

// DBS - Depth Below Surface
// https://gpsd.gitlab.io/gpsd/NMEA.html#_dbs_depth_below_surface
type DBS struct {
	BaseSentence
	DepthFeet    Float64
	DepthMeters  Float64
	DepthFathoms Float64
}

// newDBS constructor
func newDBS(s BaseSentence) (DBS, error) {
	p := NewParser(s)
	p.AssertType(TypeDBS)
	return DBS{
		BaseSentence: s,
		DepthFeet:    p.Float64(0, "depth_feet"),
		DepthMeters:  p.Float64(2, "depth_meters"),
		DepthFathoms: p.Float64(4, "depth_fathoms"),
	}, p.Err()
}

// GetDepthBelowSurface retrieves the depth below surface from the sentence
func (s DBS) GetDepthBelowSurface() (float64, error) {
	if v, err := s.DepthMeters.GetValue(); err == nil {
		return v, nil
	}
	if v, err := s.DepthMeters.GetValue(); err == nil {
		return (unit.Length(v) * unit.Foot).Meters(), nil
	}
	if v, err := s.DepthFeet.GetValue(); err == nil {
		return (unit.Length(v) * unit.Foot).Meters(), nil
	}
	if v, err := s.DepthFathoms.GetValue(); err == nil {
		return (unit.Length(v) * unit.Fathom).Meters(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
