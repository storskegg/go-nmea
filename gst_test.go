package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("GST", func() {
	var (
		sentence Sentence
		parsed   GST
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(GST)
			} else {
				parsed = GST{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPGST,172814.0,0.006,0.023,0.020,273.6,0.023,0.020,0.031*6A"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GST struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":                                 Equal(NewTime(17, 28, 14, 0)),
					"RMSPseudorangeResiduals":              Equal(NewFloat64(0.006)),
					"ErrorEllipseSemiMajorAxis1SigmaError": Equal(NewFloat64(0.023)),
					"ErrorEllipseSemiMinorAxis1SigmaError": Equal(NewFloat64(0.020)),
					"ErrorEllipseOrientation":              Equal(NewFloat64(273.6)),
					"Latitude1SigmaError":                  Equal(NewFloat64(0.023)),
					"Longitude1SigmaError":                 Equal(NewFloat64(0.020)),
					"Height1SigmaError":                    Equal(NewFloat64(0.031)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPGST,172814.0,0.006,0.023,0.020,273.6,0.023,0.020,0.031*28"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [6A != 28]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
})
