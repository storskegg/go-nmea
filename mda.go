package nmea

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

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

// GetTrueWindDirection retrieves the true wind direction from the sentence
func (s MDA) GetTrueWindDirection() (float64, error) {
	if v, err := s.WindDirectionTrue.GetValue(); err == nil {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetMagneticWindDirection retrieves the true wind direction from the sentence
func (s MDA) GetMagneticWindDirection() (float64, error) {
	if v, err := s.WindDirectionMagnetic.GetValue(); err == nil {
		return (unit.Angle(v) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetWindSpeed retrieves wind speed from the sentence
func (s MDA) GetWindSpeed() (float64, error) {
	if v, err := s.WindSpeedInMetersPerSecond.GetValue(); err == nil {
		return v, nil
	}
	if v, err := s.WindSpeedInKnots.GetValue(); err == nil {
		return (unit.Speed(v) * unit.Knot).MetersPerSecond(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetOutsideTemperature retrieves the outside air temperature from the sentence
func (s MDA) GetOutsideTemperature() (float64, error) {
	if v, err := s.AirTemperature.GetValue(); err == nil {
		return unit.FromCelsius(v).Kelvin(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetOutsideTemperature retrieves the outside air temperature from the sentence
func (s MDA) GetWaterTemperature() (float64, error) {
	if v, err := s.WaterTemperature.GetValue(); err == nil {
		return unit.FromCelsius(v).Kelvin(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetDewPointTemperature retrieves the dew point temperature from the sentence
func (s MDA) GetDewPointTemperature() (float64, error) {
	if v, err := s.DewPoint.GetValue(); err == nil {
		return unit.FromCelsius(v).Kelvin(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetOutsidePressure retrieves the outside pressure from the sentence
func (s MDA) GetOutsidePressure() (float64, error) {
	if v, err := s.BarometricPressureInBar.GetValue(); err == nil {
		return (unit.Pressure(v) * unit.Bar).Pascals(), nil
	}
	if v, err := s.BarometricPressureInInchesOfMercury.GetValue(); err == nil {
		return (unit.Pressure(v) * unit.InchOfMercury).Pascals(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetHumidity retrieves the relative humidity from the sentence
func (s MDA) GetHumidity() (float64, error) {
	if v, err := s.RelativeHumidity.GetValue(); err == nil {
		return v / 100.0, nil
	}
	return 0, fmt.Errorf("value is unavailable")
}
