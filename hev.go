package nmea

import "fmt"

const (
	// TypeHEV type for HEV sentences
	TypeHEV = "HEV"
)

// HEV - Heave
type HEV struct {
	BaseSentence
	Heave Float64 // Heave in meters
}

// newHEV constructor
func newHEV(s BaseSentence) (HEV, error) {
	p := NewParser(s)
	p.AssertType(TypeHEV)
	m := HEV{
		BaseSentence: s,
		Heave:        p.Float64(0, "heave"),
	}
	return m, p.Err()
}

// GetHeave retrieves the heave from the sentence
func (s HEV) GetHeave() (float64, error) {
	if v, err := s.Heave.GetValue(); err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
