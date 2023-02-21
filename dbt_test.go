package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("DBT", func() {
	var (
		sentence Sentence
		parsed   DBT
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(DBT)
			} else {
				parsed = DBT{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$IIDBT,032.93,f,010.04,M,005.42,F*2C"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid DBT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"DepthFeet":    Equal(NewFloat64(32.93)),
					"DepthMeters":  Equal(NewFloat64(10.04)),
					"DepthFathoms": Equal(NewFloat64(5.42)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$IIDBT,032.93,f,010.04,M,005.42,F*22"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [2C != 22]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a DBT struct", func() {
		BeforeEach(func() {
			parsed = DBT{
				DepthFeet:    NewFloat64(DepthBelowSurfaceFeet - DepthTransducerFeet),
				DepthMeters:  NewFloat64(DepthBelowSurfaceMeters - DepthTransducerMeters),
				DepthFathoms: NewFloat64(DepthBelowSurfaceFathoms - DepthTransducerFathoms),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("when having a struct with only depth in feet set", func() {
			JustBeforeEach(func() {
				parsed.DepthMeters = NewInvalidFloat64("")
				parsed.DepthFathoms = NewInvalidFloat64("")
			})
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("when having a struct with only depth in fathoms set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = NewInvalidFloat64("")
				parsed.DepthMeters = NewInvalidFloat64("")
			})
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("when having a struct with only depth in meters set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = NewInvalidFloat64("")
				parsed.DepthFathoms = NewInvalidFloat64("")
			})
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("when having a struct with missing depth values", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = NewInvalidFloat64("")
				parsed.DepthMeters = NewInvalidFloat64("")
				parsed.DepthFathoms = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowTransducer()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
