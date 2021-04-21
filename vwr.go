package nmea

const (
	TypeVWR = "VWR"

	LeftOfBow  = "L"
	RightOfBow = "R"
)

// Sentence info:
// 1    Measured angle relative to the vessel, left/right of vessel heading, to the nearest 0.1 degree
// 2    L: Left, or R: Right
// 3    Measured wind speed, knots, to the nearest 0.1 knot
// 4    N: knots
// 5    Wind speed, meters per second, to the nearest 0.1 m/s
// 6    M: meters per second
// 7    Wind speed, km/h, to the nearest km/h
// 8    K: km/h

type VWR struct {
	BaseSentence
	Angle                        Float64
	LeftRightOfBow               string
	WindSpeedInKnots             Float64
	WindSpeedInMetersPerSecond   Float64
	WindSpeedInKilometersPerHour Float64
}

// newVWR constructor
func newVWR(s BaseSentence) (VWR, error) {
	p := NewParser(s)
	p.AssertType(TypeVWR)
	m := VWR{
		BaseSentence:                 s,
		Angle:                        p.Float64(0, "Angle"),
		LeftRightOfBow:               p.EnumString(1, "LeftRightOfBow", LeftOfBow, RightOfBow),
		WindSpeedInKnots:             p.Float64(2, "WindSpeedInKnots"),
		WindSpeedInMetersPerSecond:   p.Float64(4, "WindSpeedInMetersPerSecond"),
		WindSpeedInKilometersPerHour: p.Float64(6, "WindSpeedInKilometersPerHour"),
	}
	return m, p.Err()
}
