package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ROT", func() {
	var (
		parsed ROT
	)
	Describe("Getting data from a $__ROT sentence", func() {
		BeforeEach(func() {
			parsed = ROT{
				RateOfTurn: NewFloat64(TrueDirectionDegrees),
				Status:     ValidROT,
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid rate of turn", func() {
				Expect(parsed.GetRateOfTurn()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing rate of turn", func() {
			JustBeforeEach(func() {
				parsed.RateOfTurn = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetRateOfTurn()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When having a parsed sentence with status flag set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Status = ""
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetRateOfTurn()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
