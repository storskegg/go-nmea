package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HEV", func() {
	var (
		parsed HEV
	)
	Describe("Getting data from a $__HEV sentence", func() {
		BeforeEach(func() {
			parsed = HEV{
				Heave: NewFloat64(HeaveMeters),
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid heave", func() {
				Expect(parsed.GetHeave()).To(Float64Equal(HeaveMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing heave", func() {
			JustBeforeEach(func() {
				parsed.Heave = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetHeave()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
