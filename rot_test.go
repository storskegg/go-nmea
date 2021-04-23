package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("ROT", func() {
	var (
		sentence Sentence
		parsed   ROT
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(ROT)
			} else {
				parsed = ROT{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPROT,3.1,A*33"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid ROT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RateOfTurn": Equal(NewFloat64(3.1)),
					"Status":     Equal(NewString(ValidROT)),
				}))
			})
		})
		Context("a sentence with no rot value", func() {
			BeforeEach(func() {
				raw = "$GPROT,,V*08"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid ROT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RateOfTurn": Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"Status":     Equal(NewString(InvalidROT)),
				}))
			})
		})
		Context("a sentence with status set to invalid", func() {
			BeforeEach(func() {
				raw = "$GPROT,3.1,V*24"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid ROT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RateOfTurn": Equal(NewFloat64(3.1)),
					"Status":     Equal(NewString(InvalidROT)),
				}))
			})
		})
		Context("a sentence with non existing status", func() {
			BeforeEach(func() {
				raw = "$GPROT,3.1,K*39"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid ROT struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RateOfTurn": Equal(NewFloat64(3.1)),
					"Status":     Equal(NewInvalidString("not a valid option")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPROT,3.1,A*00"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [33 != 00]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a ROT struct", func() {
		BeforeEach(func() {
			parsed = ROT{
				RateOfTurn: NewFloat64(TrueDirectionDegrees),
				Status:     NewString(ValidROT),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid rate of turn", func() {
				Expect(parsed.GetRateOfTurn()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
		})
		Context("when having a struct with missing rate of turn", func() {
			JustBeforeEach(func() {
				parsed.RateOfTurn = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetRateOfTurn()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a struct with status flag set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Status = NewInvalidString("")
			})
			It("returns an error", func() {
				_, err := parsed.GetRateOfTurn()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
