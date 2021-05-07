package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("HDT", func() {
	var (
		sentence Sentence
		parsed   HDT
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(HDT)
			} else {
				parsed = HDT{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPHDT,123.456,T*32"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid HDT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewFloat64(123.456)),
					"True":    Equal(true),
				}))
			})
		})
		Context("a sentence with an non existing true", func() {
			BeforeEach(func() {
				raw = "$GPHDT,123.456,X*3E"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid HDT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewFloat64(123.456)),
					"True":    Equal(false),
				}))
			})
		})
		Context("a sentence with an invalid heading", func() {
			BeforeEach(func() {
				raw = "$GPHDT,XXX,T*43"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid HDT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"XXX\": invalid syntax")),
					"True":    Equal(true),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPHDT,123.456,T*D7"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [32 != D7]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a HDT struct", func() {
		BeforeEach(func() {
			parsed = HDT{
				Heading: NewFloat64(TrueDirectionDegrees),
				True:    true,
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
		})
		Context("when having a struct with missing heading", func() {
			JustBeforeEach(func() {
				parsed.Heading = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a struct with true flag not set", func() {
			JustBeforeEach(func() {
				parsed.True = false
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
