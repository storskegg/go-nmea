package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("VTG", func() {
	var (
		sentence Sentence
		parsed   VTG
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(VTG)
			} else {
				parsed = VTG{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPVTG,45.5,T,67.5,M,30.45,N,56.40,K*4B"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid VTG struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TrueTrack":        Equal(NewFloat64(45.5)),
					"MagneticTrack":    Equal(NewFloat64(67.5)),
					"GroundSpeedKnots": Equal(NewFloat64(30.45)),
					"GroundSpeedKPH":   Equal(NewFloat64(56.40)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPVTG,45.5,T,67.5,M,30.45,N,56.40,K*1C"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [4B != 1C]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a VTG struct", func() {
		BeforeEach(func() {
			parsed = VTG{
				TrueTrack:        NewFloat64(TrueDirectionDegrees),
				MagneticTrack:    NewFloat64(MagneticDirectionDegrees),
				GroundSpeedKPH:   NewFloat64(SpeedOverGroundKPH),
				GroundSpeedKnots: NewFloat64(SpeedOverGroundKnots),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(BeNumerically("~", MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing true track", func() {
			JustBeforeEach(func() {
				parsed.TrueTrack = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueCourseOverGround()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(BeNumerically("~", MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing magnetic track", func() {
			JustBeforeEach(func() {
				parsed.MagneticTrack = NewInvalidFloat64("")
			})
			It("returns a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetMagneticCourseOverGround()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing speed over ground kph", func() {
			JustBeforeEach(func() {
				parsed.GroundSpeedKPH = NewInvalidFloat64("")
			})
			It("returns a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(BeNumerically("~", MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing speed over ground knots", func() {
			JustBeforeEach(func() {
				parsed.GroundSpeedKnots = NewInvalidFloat64("")
			})
			It("returns a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(BeNumerically("~", MagneticDirectionRadians, 0.00001))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("when having a struct with missing speed over ground kph and knots", func() {
			JustBeforeEach(func() {
				parsed.GroundSpeedKPH = NewInvalidFloat64("")
				parsed.GroundSpeedKnots = NewInvalidFloat64("")
			})
			It("returns a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(BeNumerically("~", MagneticDirectionRadians, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetSpeedOverGround()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
