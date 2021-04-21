package nmea

const (
	// TypeHEV type for HEV sentences
	TypeHEV = "HEV"
)

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
