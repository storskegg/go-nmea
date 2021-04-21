package nmea

import "fmt"

const (
	// TypeDPT type for DPT sentences
	TypeDPT = "DPT"
)

// DPT - Depth of Water
// https://gpsd.gitlab.io/gpsd/NMEA.html#_dpt_depth_of_water
type DPT struct {
	BaseSentence
	Depth      Float64
	Offset     Float64
	RangeScale Float64
}

// newDPT constructor
func newDPT(s BaseSentence) (DPT, error) {
	p := NewParser(s)
	p.AssertType(TypeDPT)
	return DPT{
		BaseSentence: s,
		Depth:        p.Float64(0, "depth"),
		Offset:       p.Float64(1, "offset"),
		RangeScale:   p.Float64(2, "range scale"),
	}, p.Err()
}

// GetDepthBelowTransducer retrieves the depth below the keel from the sentence
func (s DPT) GetDepthBelowTransducer() (float64, error) {
	if v, err := s.Depth.GetValue(); err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetDepthBelowKeel retrieves the depth below the keel from the sentence
func (s DPT) GetDepthBelowKeel() (float64, error) {
	if vDepth, err := s.Depth.GetValue(); err == nil {
		if vOffset, err := s.Offset.GetValue(); err == nil && vOffset < 0 {
			return vDepth + vOffset, nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetDepthBelowSurface retrieves the depth below surface from the sentence
func (s DPT) GetDepthBelowSurface() (float64, error) {
	if vDepth, err := s.Depth.GetValue(); err == nil {
		if vOffset, err := s.Offset.GetValue(); err == nil && vOffset > 0 {
			return vDepth + vOffset, nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}
