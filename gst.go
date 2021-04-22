package nmea

const (
	TypeGST = "GST"
)

// Sentence info:
// 1 	UTC of position fix
// 2 	RMS value of the pseudorange residuals; includes carrier phase residuals during periods of RTK (float) and RTK (fixed) processing
// 3 	Error ellipse semi-major axis 1 sigma error, in meters
// 4 	Error ellipse semi-minor axis 1 sigma error, in meters
// 5 	Error ellipse orientation, degrees from true north
// 6 	Latitude 1 sigma error, in meters
// 7 	Longitude 1 sigma error, in meters
// 8 	Height 1 sigma error, in meters

type GST struct {
	BaseSentence
	Time                                 Time
	RMSPseudorangeResiduals              Float64
	ErrorEllipseSemiMajorAxis1SigmaError Float64
	ErrorEllipseSemiMinorAxis1SigmaError Float64
	ErrorEllipseOrientation              Float64
	Latitude1SigmaError                  Float64
	Longitude1SigmaError                 Float64
	Height1SigmaError                    Float64
}

// newGST constructor
func newGST(s BaseSentence) (GST, error) {
	p := NewParser(s)
	p.AssertType(TypeGST)
	m := GST{
		BaseSentence:                         s,
		Time:                                 p.Time(0, "time"),
		RMSPseudorangeResiduals:              p.Float64(1, "RMSPseudorangeResiduals"),
		ErrorEllipseSemiMajorAxis1SigmaError: p.Float64(2, "ErrorEllipseSemiMajorAxis1SigmaError"),
		ErrorEllipseSemiMinorAxis1SigmaError: p.Float64(3, "ErrorEllipseSemiMinorAxis1SigmaError"),
		ErrorEllipseOrientation:              p.Float64(4, "ErrorEllipseOrientation"),
		Latitude1SigmaError:                  p.Float64(5, "Latitude1SigmaError"),
		Longitude1SigmaError:                 p.Float64(6, "Longitude1SigmaError"),
		Height1SigmaError:                    p.Float64(7, "Height1SigmaError"),
	}
	return m, p.Err()
}
