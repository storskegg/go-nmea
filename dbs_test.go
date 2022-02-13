package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("DBS", func() {
	var (
		sentence Sentence
		parsed   DBS
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(DBS)
			} else {
				parsed = DBS{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*21"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid DBS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"DepthFeet":    Equal(NewFloat64(1.9)),
					"DepthMeters":  Equal(NewFloat64(0.58)),
					"DepthFathoms": Equal(NewFloat64(0.3)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*25"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [21 != 25]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a DBS struct", func() {
		BeforeEach(func() {
			parsed = DBS{
				DepthFeet:    NewFloat64(DepthBelowSurfaceFeet),
				DepthMeters:  NewFloat64(DepthBelowSurfaceMeters),
				DepthFathoms: NewFloat64(DepthBelowSurfaceFathoms),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(BeNumerically("~", DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("when having a struct with only depth in feet set", func() {
			JustBeforeEach(func() {
				parsed.DepthMeters = NewInvalidFloat64("")
				parsed.DepthFathoms = NewInvalidFloat64("")
			})
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(BeNumerically("~", DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("when having a struct with only depth in fathoms set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = NewInvalidFloat64("")
				parsed.DepthMeters = NewInvalidFloat64("")
			})
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(BeNumerically("~", DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("when having a struct with only depth in meters set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = NewInvalidFloat64("")
				parsed.DepthFathoms = NewInvalidFloat64("")
			})
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(BeNumerically("~", DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("when having a struct with missing depth values", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = NewInvalidFloat64("")
				parsed.DepthMeters = NewInvalidFloat64("")
				parsed.DepthFathoms = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
