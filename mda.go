package nmea

const (
	TypeMDA = "MDA"
)

// Sentence info:
// 1    Barometric pressure, inches of mercury, to the nearest 0.01 inch
// 2    I = inches of mercury
// 3    Barometric pressure, bars, to the nearest .001 bar
// 4    B = bars
// 5    Air temperature, degrees C, to the nearest 0.1 degree C
// 6    C = degrees C
// 7    Water temperature, degrees C (this field left blank by WeatherStation)
// 8    C = degrees C (this field left blank by WeatherStation)
// 9    Relative humidity, percent, to the nearest 0.1 percent
// 10   Absolute humidity, percent (this field left blank by WeatherStation)
// 11   Dew point, degrees C, to the nearest 0.1 degree C
// 12   C = degrees C
// 13   Wind direction, degrees True, to the nearest 0.1 degree
// 14   T = true
// 15   Wind direction, degrees Magnetic, to the nearest 0.1 degree
// 16   M = magnetic
// 17   Wind speed, knots, to the nearest 0.1 knot
// 18   N = knots
// 19   Wind speed, meters per second, to the nearest 0.1 m/s
// 20   M = meters per second

type MDA struct {
	BaseSentence
	BarometricPressureInInchesOfMercury Float64
	BarometricPressureInBar             Float64
	AirTemperature                      Float64
	WaterTemperature                    Float64
	RelativeHumidity                    Float64
	DewPoint                            Float64
	WindDirectionTrue                   Float64
	WindDirectionMagnetic               Float64
	WindSpeedInKnots                    Float64
	WindSpeedInMetersPerSecond          Float64
}

// newMDA constructor
func newMDA(s BaseSentence) (MDA, error) {
	p := NewParser(s)
	p.AssertType(TypeMDA)
	m := MDA{
		BaseSentence:                        s,
		BarometricPressureInInchesOfMercury: p.Float64(0, "BarometricPressureInInchesOfMercury"),
		BarometricPressureInBar:             p.Float64(2, "BarometricPressureInBar"),
		AirTemperature:                      p.Float64(4, "AirTemperature"),
		WaterTemperature:                    p.Float64(6, "WaterTemperature"),
		RelativeHumidity:                    p.Float64(8, "RelativeHumidity"),
		DewPoint:                            p.Float64(10, "DewPoint"),
		WindDirectionTrue:                   p.Float64(12, "WindDirectionTrue"),
		WindDirectionMagnetic:               p.Float64(14, "WindDirectionMagnetic"),
		WindSpeedInKnots:                    p.Float64(16, "WindSpeedInKnots"),
		WindSpeedInMetersPerSecond:          p.Float64(18, "WindSpeedInMetersPerSecond"),
	}
	return m, p.Err()
}
