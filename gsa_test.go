package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("GSA", func() {
	var (
		sentence Sentence
		parsed   GSA
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(GSA)
			} else {
				parsed = GSA{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPGSA,A,3,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*36"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Mode":    Equal(NewString(Auto)),
					"FixType": Equal(NewString(Fix3D)),
					"SV":      Equal([]String{NewString("22"), NewString("19"), NewString("18"), NewString("27"), NewString("14"), NewString("03")}),
					"PDOP":    Equal(NewFloat64(3.1)),
					"HDOP":    Equal(NewFloat64(2)),
					"VDOP":    Equal(NewFloat64(2.4)),
				}))
			})
		})
		Context("a sentence with a bad mode", func() {
			BeforeEach(func() {
				raw = "$GPGSA,F,3,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*31"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Mode":    Equal(NewInvalidString("not a valid option")),
					"FixType": Equal(NewString(Fix3D)),
					"SV":      Equal([]String{NewString("22"), NewString("19"), NewString("18"), NewString("27"), NewString("14"), NewString("03")}),
					"PDOP":    Equal(NewFloat64(3.1)),
					"HDOP":    Equal(NewFloat64(2)),
					"VDOP":    Equal(NewFloat64(2.4)),
				}))
			})
		})
		Context("a sentence with a bad fix", func() {
			BeforeEach(func() {
				raw = "$GPGSA,A,6,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*33"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Mode":    Equal(NewString(Auto)),
					"FixType": Equal(NewInvalidString("not a valid option")),
					"SV":      Equal([]String{NewString("22"), NewString("19"), NewString("18"), NewString("27"), NewString("14"), NewString("03")}),
					"PDOP":    Equal(NewFloat64(3.1)),
					"HDOP":    Equal(NewFloat64(2)),
					"VDOP":    Equal(NewFloat64(2.4)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPGSA,A,3,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*00"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [36 != 00]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a GSA struct", func() {
		BeforeEach(func() {
			parsed = GSA{
				Mode:    NewString(Auto),
				FixType: NewString(Fix3D),
				SV:      make([]String, Satellites),
				PDOP:    NewInvalidFloat64(""),
				HDOP:    NewInvalidFloat64(""),
				VDOP:    NewInvalidFloat64(""),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("returns a valid fix type", func() {
				Expect(parsed.GetFixType()).To(Equal(Fix3D))
			})
		})
	})
})
