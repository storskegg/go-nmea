package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("DPT", func() {
	var (
		sentence Sentence
		parsed   DPT
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(DPT)
			} else {
				parsed = DPT{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$SDDPT,0.5,0.5,*7B"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid DPT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Depth":      Equal(NewFloat64(0.5)),
					"Offset":     Equal(NewFloat64(0.5)),
					"RangeScale": Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
				}))
			})
		})
		Context("A valid sentence with range scale", func() {
			BeforeEach(func() {
				raw = "$SDDPT,0.5,0.5,0.1*54"
			})
			Specify("no error is returned", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid DPT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Depth":      Equal(NewFloat64(0.5)),
					"Offset":     Equal(NewFloat64(0.5)),
					"RangeScale": Equal(NewFloat64(0.1)),
				}))
			})
		})
		Context("Bad checksum", func() {
			BeforeEach(func() {
				raw = "$SDDPT,0.5,0.5,*AA"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [7B != AA]"))
			})
		})
	})
	Describe("Getting data from a DPT struct", func() {
		BeforeEach(func() {
			parsed = DPT{
				Depth: NewFloat64(DepthBelowSurfaceMeters - DepthTransducerMeters),
			}
		})
		Context("When having a parsed sentence and a positive offset", func() {
			JustBeforeEach(func() {
				parsed.Offset = NewFloat64(DepthTransducerMeters)
			})
			It("returns a valid depth below transducer", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
			It("returns a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(BeNumerically("~", DepthBelowSurfaceMeters, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowKeel()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When having a parsed sentence and a negative offset", func() {
			JustBeforeEach(func() {
				parsed.Offset = NewFloat64(DepthTransducerMeters - DepthKeelMeters)
			})
			It("returns a valid depth below transducer", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid depth below keel", func() {
				Expect(parsed.GetDepthBelowKeel()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthKeelMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence and no offset", func() {
			JustBeforeEach(func() {
				parsed.Offset = NewInvalidFloat64("")
			})
			It("returns a valid depth below transducer", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(BeNumerically("~", DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowKeel()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When having a parsed sentence and no depth", func() {
			JustBeforeEach(func() {
				parsed.Depth = NewInvalidFloat64("")
				parsed.Offset = NewFloat64(DepthTransducerMeters)
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowTransducer()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetDepthBelowKeel()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
