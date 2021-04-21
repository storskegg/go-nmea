package nmea

const (
	TypeMWV = "MWV"

	ReferenceRelative = "R"
	ReferenceTrue     = "T"

	WindSpeedUnitKPH   = "K"
	WindSpeedUnitKnots = "N"
	WindSpeedUnitMPS   = "M"
	WindSpeedUnitMPH   = "S"

	ValidMWV   = "A"
	InvalidMWV = "V"
)

// Sentence info:
// 1    Wind angle, 0.0 to 359.9 degrees, in relation to the vessel’s bow/centerline, to the nearest 0.1 degree. If the data for this field is not valid, the field will be blank.
// 2    Reference:
// 			R: Relative (apparent wind, as felt when standing on the moving ship)
// 			T: Theoretical (calculated actual wind, as though the vessel were stationary)
// 3    Wind speed, to the nearest tenth of a unit.  If the data for this field is not valid, the field will be blank.
// 4    Wind speed units:
//			K: km/hr
//			M: m/s
//			N: knots
//			S: statute miles/hr
// 5    Status:
//			A: data valid
//			V: data invalid

type MWV struct {
	BaseSentence
	Angle         Float64
	Reference     string
	WindSpeed     Float64
	WindSpeedUnit string
	Status        string
}

// newMWV constructor
func newMWV(s BaseSentence) (MWV, error) {
	p := NewParser(s)
	p.AssertType(TypeMWV)
	m := MWV{
		BaseSentence:  s,
		Angle:         p.Float64(0, "Angle"),
		Reference:     p.EnumString(1, "Reference", ReferenceRelative, ReferenceTrue),
		WindSpeed:     p.Float64(2, "WindSpeed"),
		WindSpeedUnit: p.EnumString(3, "WindSpeedUnit", WindSpeedUnitKPH, WindSpeedUnitKnots, WindSpeedUnitMPS, WindSpeedUnitMPH),
		Status:        p.EnumString(4, "Status", ValidMWV, InvalidMWV),
	}
	return m, p.Err()
}
