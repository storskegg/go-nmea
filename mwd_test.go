package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MWD", func() {
	var (
		parsed MWD
	)
	Describe("Getting data from a $__MWD sentence", func() {
		BeforeEach(func() {
			parsed = MWD{
				WindDirectionTrue:          NewFloat64(TrueDirectionDegrees),
				WindDirectionMagnetic:      NewFloat64(MagneticDirectionDegrees),
				WindSpeedInKnots:           NewFloat64(SpeedOverGroundKnots),
				WindSpeedInMetersPerSecond: NewFloat64(SpeedOverGroundMPS),
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid true wind direction", func() {
				Expect(parsed.GetTrueWindDirection()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic wind direction", func() {
				Expect(parsed.GetMagneticWindDirection()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing data", func() {
			JustBeforeEach(func() {
				parsed = MWD{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetMagneticWindDirection()
				Expect(err).To(HaveOccurred())
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetTrueWindDirection()
				Expect(err).To(HaveOccurred())
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetWindSpeed()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
