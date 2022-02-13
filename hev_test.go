package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("HEV", func() {
	var (
		sentence Sentence
		parsed   HEV
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(HEV)
			} else {
				parsed = HEV{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPHEV,-0.07*54"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid HEV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heave": Equal(NewFloat64(-0.07)),
				}))
			})
		})
		Context("a sentence with an invalid heave", func() {
			BeforeEach(func() {
				raw = "$GPHEV,-B.07*26"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid HEV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heave": Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"-B.07\": invalid syntax")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPHEV,-0.07*00"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [54 != 00]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a HEV struct", func() {
		BeforeEach(func() {
			parsed = HEV{
				Heave: NewFloat64(HeaveMeters),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid heave", func() {
				Expect(parsed.GetHeave()).To(BeNumerically("~", HeaveMeters, 0.00001))
			})
		})
		Context("when having a struct with missing heave", func() {
			JustBeforeEach(func() {
				parsed.Heave = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetHeave()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
