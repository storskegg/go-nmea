package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("MDA", func() {
	var (
		sentence Sentence
		parsed   MDA
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(MDA)
			} else {
				parsed = MDA{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$WIMDA,30.1176,I,1.0199,B,44.0,C,12.7,C,78.9,,14.2,C,359.0,T,358.7,M,6.4,N,3.3,M*37"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MDA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"BarometricPressureInInchesOfMercury": Equal(NewFloat64(30.1176)),
					"BarometricPressureInBar":             Equal(NewFloat64(1.0199)),
					"AirTemperature":                      Equal(NewFloat64(44)),
					"WaterTemperature":                    Equal(NewFloat64(12.7)),
					"RelativeHumidity":                    Equal(NewFloat64(78.9)),
					"DewPoint":                            Equal(NewFloat64(14.2)),
					"WindDirectionTrue":                   Equal(NewFloat64(359)),
					"WindDirectionMagnetic":               Equal(NewFloat64(358.7)),
					"WindSpeedInKnots":                    Equal(NewFloat64(6.4)),
					"WindSpeedInMetersPerSecond":          Equal(NewFloat64(3.3)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$WIMDA,30.1383,I,1.0206,B,43.9,C,,,,,,,7.4,T,7.1,M,7.7,N,4.0,M*2B"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [2A != 2B]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a MDA struct", func() {
		BeforeEach(func() {
			parsed = MDA{
				BarometricPressureInInchesOfMercury: NewFloat64(PressureInchesOfMercury),
				BarometricPressureInBar:             NewFloat64(PressureBar),
				AirTemperature:                      NewFloat64(AirTemperatureCelcius),
				WaterTemperature:                    NewFloat64(WaterTemperatureCelcius),
				RelativeHumidity:                    NewFloat64(RelativeHumidityPercentage),
				DewPoint:                            NewFloat64(DewPointCelcius),
				WindDirectionTrue:                   NewFloat64(TrueDirectionDegrees),
				WindDirectionMagnetic:               NewFloat64(MagneticDirectionDegrees),
				WindSpeedInKnots:                    NewFloat64(SpeedOverGroundKnots),
				WindSpeedInMetersPerSecond:          NewFloat64(SpeedOverGroundMPS),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid outside pressure", func() {
				Expect(parsed.GetOutsidePressure()).To(BeNumerically("~", PressurePascal, 0.00001))
			})
			It("returns a valid outside temperature", func() {
				Expect(parsed.GetOutsideTemperature()).To(BeNumerically("~", AirTemperatureKelvin, 0.00001))
			})
			It("returns a valid water temperature", func() {
				Expect(parsed.GetWaterTemperature()).To(BeNumerically("~", WaterTemperatureKelvin, 0.00001))
			})
			It("returns a valid humidity", func() {
				Expect(parsed.GetHumidity()).To(BeNumerically("~", RelativeHumidityRatio, 0.00001))
			})
			It("returns a valid dew point temperature", func() {
				Expect(parsed.GetDewPointTemperature()).To(BeNumerically("~", DewPointKelvin, 0.00001))
			})
			It("returns a valid true wind direction", func() {
				Expect(parsed.GetTrueWindDirection()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic wind direction", func() {
				Expect(parsed.GetMagneticWindDirection()).To(BeNumerically("~", MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with pressure inches of mercury missing", func() {
			JustBeforeEach(func() {
				parsed.BarometricPressureInInchesOfMercury = NewInvalidFloat64("")
			})
			It("returns a valid outside pressure", func() {
				Expect(parsed.GetOutsidePressure()).To(BeNumerically("~", PressurePascal, 0.00001))
			})
		})
		Context("when having a struct with pressure bar missing", func() {
			JustBeforeEach(func() {
				parsed.BarometricPressureInBar = NewInvalidFloat64("")
			})
			It("returns a valid outside pressure", func() {
				Expect(parsed.GetOutsidePressure()).To(BeNumerically("~", PressurePascal, 0.5))
			})
		})
		Context("when having a struct with wind speed in meters per second missing", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInMetersPerSecond = NewInvalidFloat64("")
			})
			It("returns a valid outside pressure", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing data", func() {
			JustBeforeEach(func() {
				parsed = MDA{}
			})
			It("returns an error", func() {
				_, err := parsed.GetOutsidePressure()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetOutsideTemperature()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetWaterTemperature()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetHumidity()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetDewPointTemperature()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetMagneticWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetWindSpeed()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
