package nmea

const (
	// TypeRMC type for RMC sentences
	TypeRMC = "RMC"
	// ValidRMC character
	ValidRMC = "A"
	// InvalidRMC character
	InvalidRMC = "V"
)

// RMC is the Recommended Minimum Specific GNSS data.
// http://aprs.gids.nl/nmea/#rmc
type RMC struct {
	BaseSentence
	Time      Time    // Time Stamp
	Validity  string  // validity - A-ok, V-invalid
	Latitude  Float64 // Latitude
	Longitude Float64 // Longitude
	Speed     Float64 // Speed in knots
	Course    Float64 // True course
	Date      Date    // Date
	Variation Float64 // Magnetic variation
}

// newRMC constructor
func newRMC(s BaseSentence) (RMC, error) {
	p := NewParser(s)
	p.AssertType(TypeRMC)
	m := RMC{
		BaseSentence: s,
		Time:         p.Time(0, "time"),
		Validity:     p.EnumString(1, "validity", ValidRMC, InvalidRMC),
		Latitude:     p.LatLong(2, 3, "latitude"),
		Longitude:    p.LatLong(4, 5, "longitude"),
		Speed:        p.Float64(6, "speed"),
		Course:       p.Float64(7, "course"),
		Date:         p.Date(8, "date"),
		Variation:    p.Float64(9, "variation"),
	}
	if m.Variation.Valid && p.EnumString(10, "direction", West, East) == West {
		m.Variation.Value = 0 - m.Variation.Value
	}
	return m, p.Err()
}
