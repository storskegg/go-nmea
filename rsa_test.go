package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("RSA", func() {
	var (
		sentence Sentence
		parsed   RSA
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(RSA)
			} else {
				parsed = RSA{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$RIRSA,3.1,A,-3.1,A*76"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RSA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RudderAngleStarboard": Equal(NewFloat64(3.1)),
					"StatusStarboard":      Equal(NewString(ValidRSA)),
					"RudderAnglePortside":  Equal(NewFloat64(-3.1)),
					"StatusPortside":       Equal(NewString(ValidRSA)),
				}))
			})
		})
		Context("a sentence with no rsa values", func() {
			BeforeEach(func() {
				raw = "$RIRSA,,V,,V*5B"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RSA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RudderAngleStarboard": Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"StatusStarboard":      Equal(NewString(InvalidRSA)),
					"RudderAnglePortside":  Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"StatusPortside":       Equal(NewString(InvalidRSA)),
				}))
			})
		})
		Context("a sentence with status set to invalid", func() {
			BeforeEach(func() {
				raw = "$RIRSA,3.1,V,-3.1,V*76"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RSA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RudderAngleStarboard": Equal(NewFloat64(3.1)),
					"StatusStarboard":      Equal(NewString(InvalidRSA)),
					"RudderAnglePortside":  Equal(NewFloat64(-3.1)),
					"StatusPortside":       Equal(NewString(InvalidRSA)),
				}))
			})
		})
		Context("a sentence with non existing status", func() {
			BeforeEach(func() {
				raw = "$RIRSA,3.1,K,-3.1,K*76"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RSA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"RudderAngleStarboard": Equal(NewFloat64(3.1)),
					"StatusStarboard":      Equal(NewInvalidString("not a valid option")),
					"RudderAnglePortside":  Equal(NewFloat64(-3.1)),
					"StatusPortside":       Equal(NewInvalidString("not a valid option")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$RIRSA,3.1,A,-3.1,A*75"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [76 != 75]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a RSA struct", func() {
		BeforeEach(func() {
			parsed = RSA{
				RudderAngleStarboard: NewFloat64(TrueDirectionDegrees),
				StatusStarboard:      NewString(ValidRSA),
				RudderAnglePortside:  NewFloat64(-TrueDirectionDegrees),
				StatusPortside:       NewString(ValidRSA),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid rate of turn", func() {
				Expect(parsed.GetRudderAngleStarboard()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
				Expect(parsed.GetRudderAnglePortside()).To(BeNumerically("~", -TrueDirectionRadians, 0.00001))
				_, err := parsed.GetRudderAngle()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a struct with a missing portside rudder angle", func() {
			JustBeforeEach(func() {
				parsed.RudderAnglePortside = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetRudderAnglePortside()
				Expect(err).To(HaveOccurred())
				Expect(parsed.GetRudderAngleStarboard()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
				Expect(parsed.GetRudderAngle()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
		})
		Context("when having a struct with status flag set to invalid", func() {
			JustBeforeEach(func() {
				parsed.StatusStarboard = NewInvalidString("")
			})
			It("returns an error", func() {
				_, err := parsed.GetRudderAngleStarboard()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
