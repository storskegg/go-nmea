package nmea

const (
	TypeMWD = "MWD"
)

// Sentence info:
// 1    Wind direction, 0.0 to 359.9 degrees True, to the nearest 0.1 degree
// 2    T: True
// 3    Wind direction, 0.0 to 359.9 degrees Magnetic, to the nearest 0.1 degree
// 4    M: Magnetic
// 5    Wind speed, knots, to the nearest 0.1 knot.
// 6    N: Knots
// 7    Wind speed, meters/second, to the nearest 0.1 m/s.
// 8    M: Meters/second

type MWD struct {
	BaseSentence
	WindDirectionTrue          Float64
	WindDirectionMagnetic      Float64
	WindSpeedInKnots           Float64
	WindSpeedInMetersPerSecond Float64
}

// newMWD constructor
func newMWD(s BaseSentence) (MWD, error) {
	p := NewParser(s)
	p.AssertType(TypeMWD)
	m := MWD{
		BaseSentence:               s,
		WindDirectionTrue:          p.Float64(0, "WindDirectionTrue"),
		WindDirectionMagnetic:      p.Float64(2, "WindDirectionMagnetic"),
		WindSpeedInKnots:           p.Float64(4, "WindSpeedInKnots"),
		WindSpeedInMetersPerSecond: p.Float64(6, "WindSpeedInMetersPerSecond"),
	}
	return m, p.Err()
}
