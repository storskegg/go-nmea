package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VWR", func() {
	var (
		parsed VWR
	)
	Describe("Getting data from a VWR struct", func() {
		BeforeEach(func() {
			parsed = VWR{
				Angle:                        NewFloat64(RelativeDirectionDegrees),
				WindSpeedInKnots:             NewFloat64(SpeedOverGroundKnots),
				WindSpeedInMetersPerSecond:   NewFloat64(SpeedOverGroundMPS),
				WindSpeedInKilometersPerHour: NewFloat64(SpeedOverGroundKPH),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(Float64Equal(RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in meters per second", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInMetersPerSecond = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(Float64Equal(RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in kilometer per hour", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInKilometersPerHour = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(Float64Equal(RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing wind speed in knots", func() {
			JustBeforeEach(func() {
				parsed.WindSpeedInKnots = NewInvalidFloat64("")
			})
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(Float64Equal(RelativeDirectionRadians, 0.00001))
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
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
