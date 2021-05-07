package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("MWV", func() {
	var (
		sentence Sentence
		parsed   MWV
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(MWV)
			} else {
				parsed = MWV{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$WIMWV,117.5,R,4.6,N,A*23"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MWV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Angle":         Equal(NewFloat64(117.5)),
					"Reference":     Equal(NewString(ReferenceRelative)),
					"WindSpeed":     Equal(NewFloat64(4.6)),
					"WindSpeedUnit": Equal(NewString(WindSpeedUnitKnots)),
					"Status":        Equal(NewString(ValidMWV)),
				}))
			})
		})
		Context("a sentence with an invalid reference", func() {
			BeforeEach(func() {
				raw = "$WIMWV,117.5,B,4.6,N,A*33"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MWV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Angle":         Equal(NewFloat64(117.5)),
					"Reference":     Equal(NewInvalidString("not a valid option")),
					"WindSpeed":     Equal(NewFloat64(4.6)),
					"WindSpeedUnit": Equal(NewString(WindSpeedUnitKnots)),
					"Status":        Equal(NewString(ValidMWV)),
				}))
			})
		})
		Context("a sentence with an invalid wind speed unit", func() {
			BeforeEach(func() {
				raw = "$WIMWV,117.5,R,4.6,L,A*21"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MWV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Angle":         Equal(NewFloat64(117.5)),
					"Reference":     Equal(NewString(ReferenceRelative)),
					"WindSpeed":     Equal(NewFloat64(4.6)),
					"WindSpeedUnit": Equal(NewInvalidString("not a valid option")),
					"Status":        Equal(NewString(ValidMWV)),
				}))
			})
		})
		Context("a sentence with status set to invalid", func() {
			BeforeEach(func() {
				raw = "$WIMWV,117.5,R,4.6,N,V*34"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MWV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Angle":         Equal(NewFloat64(117.5)),
					"Reference":     Equal(NewString(ReferenceRelative)),
					"WindSpeed":     Equal(NewFloat64(4.6)),
					"WindSpeedUnit": Equal(NewString(WindSpeedUnitKnots)),
					"Status":        Equal(NewString(InvalidMWV)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$WIMWV,117.5,R,4.6,N,A*42"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [23 != 42]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a MWV struct", func() {
		BeforeEach(func() {
			parsed = MWV{
				Angle:         NewFloat64(RelativeDirectionDegrees),
				Reference:     NewString(ReferenceRelative),
				WindSpeed:     NewFloat64(SpeedOverGroundMPS),
				WindSpeedUnit: NewString(WindSpeedUnitMPS),
				Status:        NewString(ValidMWV),
			}
		})
		Context("when having a struct with reference set to relative", func() {
			It("returns a valid relative wind direction", func() {
				Expect(parsed.GetRelativeWindDirection()).To(BeNumerically("~", RelativeDirectionRadians, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with reference set to true", func() {
			JustBeforeEach(func() {
				parsed.Angle = NewFloat64(TrueDirectionDegrees)
				parsed.Reference = NewString(ReferenceTrue)
			})
			It("returns an error", func() {
				_, err := parsed.GetRelativeWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid true wind direction", func() {
				Expect(parsed.GetTrueWindDirection()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
		})
		Context("when having a struct with wind speed in kmh", func() {
			JustBeforeEach(func() {
				parsed.WindSpeed = NewFloat64(SpeedOverGroundKPH)
				parsed.WindSpeedUnit = NewString(WindSpeedUnitKPH)
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with wind speed in knots", func() {
			JustBeforeEach(func() {
				parsed.WindSpeed = NewFloat64(SpeedOverGroundKnots)
				parsed.WindSpeedUnit = NewString(WindSpeedUnitKnots)
			})
			It("returns a valid wind speed", func() {
				Expect(parsed.GetWindSpeed()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with an invalid wind speed unit", func() {
			JustBeforeEach(func() {
				parsed.WindSpeed = NewFloat64(SpeedOverGroundKnots)
				parsed.WindSpeedUnit = NewString("A")
			})
			It("returns an error", func() {
				_, err := parsed.GetWindSpeed()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a struct with status set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Status = NewInvalidString("")
			})
			It("returns an error", func() {
				_, err := parsed.GetRelativeWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetWindSpeed()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a struct with missing data", func() {
			JustBeforeEach(func() {
				parsed = MWV{}
			})
			It("returns an error", func() {
				_, err := parsed.GetRelativeWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueWindDirection()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetWindSpeed()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
