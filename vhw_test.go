package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("VHW", func() {
	var (
		sentence Sentence
		parsed   VHW
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(VHW)
			} else {
				parsed = VHW{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$VWVHW,45.0,T,43.0,M,3.5,N,6.4,K*56"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid VHW struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TrueHeading":            Equal(NewFloat64(45.0)),
					"MagneticHeading":        Equal(NewFloat64(43.0)),
					"SpeedThroughWaterKnots": Equal(NewFloat64(3.5)),
					"SpeedThroughWaterKPH":   Equal(NewFloat64(6.4)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$VWVHW,45.0,T,43.0,M,3.5,N,6.4,K*7F"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [56 != 7F]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a VHW struct", func() {
		BeforeEach(func() {
			parsed = VHW{
				TrueHeading:            NewFloat64(TrueDirectionDegrees),
				MagneticHeading:        NewFloat64(MagneticDirectionDegrees),
				SpeedThroughWaterKPH:   NewFloat64(SpeedThroughWaterKPH),
				SpeedThroughWaterKnots: NewFloat64(SpeedThroughWaterKnots),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("when having a struct with missing true heading", func() {
			JustBeforeEach(func() {
				parsed.TrueHeading = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("when having a struct with missing magnetic track", func() {
			JustBeforeEach(func() {
				parsed.MagneticHeading = NewInvalidFloat64("")
			})
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetMagneticHeading()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("when having a struct with missing speed over ground kph", func() {
			JustBeforeEach(func() {
				parsed.SpeedThroughWaterKPH = NewInvalidFloat64("")
			})
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("when having a struct with missing speed over ground knots", func() {
			JustBeforeEach(func() {
				parsed.SpeedThroughWaterKnots = NewInvalidFloat64("")
			})
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("when having a struct with missing speed over ground kph and knots", func() {
			JustBeforeEach(func() {
				parsed.SpeedThroughWaterKPH = NewInvalidFloat64("")
				parsed.SpeedThroughWaterKnots = NewInvalidFloat64("")
			})
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetSpeedThroughWater()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
