package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("THS", func() {
	var (
		sentence Sentence
		parsed   THS
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(THS)
			} else {
				parsed = THS{}
			}
		})
		Context("a valid sentence with status set to autonomous", func() {
			BeforeEach(func() {
				raw = "$INTHS,123.456,A*20"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid THS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewFloat64(123.456)),
					"Status":  Equal(NewString(AutonomousTHS)),
				}))
			})
		})
		Context("a valid sentence with status set to estimated", func() {
			BeforeEach(func() {
				raw = "$INTHS,123.456,E*24"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid THS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewFloat64(123.456)),
					"Status":  Equal(NewString(EstimatedTHS)),
				}))
			})
		})
		Context("a valid sentence with status set to manual", func() {
			BeforeEach(func() {
				raw = "$INTHS,123.456,M*2C"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid THS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewFloat64(123.456)),
					"Status":  Equal(NewString(ManualTHS)),
				}))
			})
		})
		Context("a valid sentence with status set to simulator", func() {
			BeforeEach(func() {
				raw = "$INTHS,123.456,S*32"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid THS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewFloat64(123.456)),
					"Status":  Equal(NewString(SimulatorTHS)),
				}))
			})
		})
		Context("a valid sentence with status set to invalid", func() {
			BeforeEach(func() {
				raw = "$INTHS,,V*1E"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid THS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"Status":  Equal(NewString(InvalidTHS)),
				}))
			})
		})
		Context("a sentence with a non existing status", func() {
			BeforeEach(func() {
				raw = "$INTHS,123.456,B*23"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid THS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewFloat64(123.456)),
					"Status":  Equal(NewInvalidString("not a valid option")),
				}))
			})
		})
		Context("a sentence with an invalid heading", func() {
			BeforeEach(func() {
				raw = "$INTHS,XXX,A*51"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid THS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Heading": Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"XXX\": invalid syntax")),
					"Status":  Equal(NewString(AutonomousTHS)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$INTHS,123.456,A*21"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [20 != 21]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a THS struct", func() {
		BeforeEach(func() {
			parsed = THS{
				Heading: NewFloat64(TrueDirectionDegrees),
				Status:  NewString(SimulatorTHS),
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
		Context("when having a struct with status flag set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Status = NewString(InvalidTHS)
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
