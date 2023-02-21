package nmea_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	MagneticDirectionRadians = float64(4.0631264986427995)
	MagneticDirectionDegrees = float64(232.8)

	TrueDirectionRadians = float64(4.094542425178697)
	TrueDirectionDegrees = float64(234.6)

	RelativeDirectionRadians = float64(-0.65973445725)
	RelativeDirectionDegrees = float64(-37.8)

	MagneticVariationRadians = float64(0.04014257)
	MagneticVariationDegrees = float64(2.3)

	SpeedOverGroundMPS   = float64(3.7222252000000005)
	SpeedOverGroundKPH   = float64(13.4)
	SpeedOverGroundKnots = float64(7.235418)

	SpeedThroughWaterMPS   = float64(3.1388889)
	SpeedThroughWaterKPH   = float64(11.3)
	SpeedThroughWaterKnots = float64(6.1015119)

	Longitude = float64(2.294481)
	Latitude  = float64(48.858372)
	Altitude  = float64(2.3)

	Satellites = int64(11)

	DepthBelowSurfaceMeters  = float64(3.9)
	DepthBelowSurfaceFeet    = float64(12.79528)
	DepthBelowSurfaceFathoms = float64(2.1325459)

	DepthTransducerMeters  = float64(1.8)
	DepthTransducerFeet    = float64(5.905512)
	DepthTransducerFathoms = float64(0.98425197)

	DepthKeelMeters  = float64(1.95)
	DepthKeelFeet    = float64(6.397638)
	DepthKeelFathoms = float64(1.0662730)

	PressurePascal          = 101600
	PressureBar             = 1.016
	PressureInchesOfMercury = 30.0026

	AirTemperatureKelvin  = 290.65
	AirTemperatureCelsius = 17.5

	WaterTemperatureKelvin  = 282.05
	WaterTemperatureCelsius = 8.9

	DewPointKelvin  = 280.35
	DewPointCelsius = 7.2

	RelativeHumidityPercentage = 55.6
	RelativeHumidityRatio      = 0.556

	HeaveMeters = -0.4
)

func TestGoNmea(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoNmea Suite")
}
