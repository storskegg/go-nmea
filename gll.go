package nmea

import "fmt"

const (
	// TypeGLL type for GLL sentences
	TypeGLL = "GLL"
	// ValidGLL character
	ValidGLL = "A"
	// InvalidGLL character
	InvalidGLL = "V"
)

// GLL is Geographic Position, Latitude / Longitude and time.
// http://aprs.gids.nl/nmea/#gll
type GLL struct {
	BaseSentence
	Latitude  Float64 // Latitude
	Longitude Float64 // Longitude
	Time      Time    // Time Stamp
	Validity  string  // validity - A-valid
}

// newGLL constructor
func newGLL(s BaseSentence) (GLL, error) {
	p := NewParser(s)
	p.AssertType(TypeGLL)
	return GLL{
		BaseSentence: s,
		Latitude:     p.LatLong(0, 1, "latitude"),
		Longitude:    p.LatLong(2, 3, "longitude"),
		Time:         p.Time(4, "time"),
		Validity:     p.EnumString(5, "validity", ValidGLL, InvalidGLL),
	}, p.Err()
}

// GetPosition2D retrieves the 2D position from the sentence
func (s GLL) GetPosition2D() (float64, float64, error) {
	if s.Validity == ValidGLL {
		if vLat, err := s.Latitude.GetValue(); err == nil {
			if vLon, err := s.Longitude.GetValue(); err == nil {
				return vLat, vLon, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("value is unavailable")
}
