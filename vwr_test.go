package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("VWR", func() {
	var (
		sentence Sentence
		parsed   VWR
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(VWR)
			} else {
				parsed = VWR{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$IIVWR,045.0,L,12.6,N,6.5,M,23.3,K*52"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid VWR struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Angle":                        Equal(NewFloat64(45)),
					"LeftRightOfBow":               Equal(NewString(LeftOfBow)),
					"WindSpeedInKnots":             Equal(NewFloat64(12.6)),
					"WindSpeedInMetersPerSecond":   Equal(NewFloat64(6.5)),
					"WindSpeedInKilometersPerHour": Equal(NewFloat64(23.3)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$IIVWR,045.0,L,12.6,N,6.5,M,23.3,K*25"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [52 != 25]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a VWR struct", func() {
		BeforeEach(func() {
			parsed = VWR{
				Angle:                        NewFloat64(RelativeDirectionDegrees),
				LeftRightOfBow:               NewString(RightOfBow),
				WindSpeedInKnots:             NewFloat64(SpeedOverGroundKnots),
				WindSpeedInMetersPerSecond:   NewFloat64(SpeedOverGroundMPS),
				WindSpeedInKilometersPerHour: NewFloat64(SpeedOverGroundKPH),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a complete struct with angle set to left of bow", func() {
			JustBeforeEach(func() {
				parsed.LeftRightOfBow = NewString(LeftOfBow)
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", 0-RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in meters per second", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInMetersPerSecond = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in kilometer per hour", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInKilometersPerHour = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in knots", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInKnots = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in meters per second and wind speed in knots", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInMetersPerSecond = NewInvalidFloat64("")
				parsed.WindSpeedInKnots = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in meters per second and wind speed in kilometer per hour", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInMetersPerSecond = NewInvalidFloat64("")
				parsed.WindSpeedInKilometersPerHour = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing data", func() {
			JustBeforeEach(func() {
				parsed = VWR{}
			})
			It("returns an error", func() {
				_, err := parsed.GetRelativeWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetWindSpeed()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
