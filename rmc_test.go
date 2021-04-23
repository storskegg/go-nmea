package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("RMC", func() {
	var (
		sentence Sentence
		parsed   RMC
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(RMC)
			} else {
				parsed = RMC{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPRMC,114509.30,A,5142.01288,N,00452.01197,E,0.0,345.6,230421,0.3,E,A,C*5C"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RMC struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":      Equal(NewTime(11, 45, 9, 300)),
					"Validity":  Equal(NewString(ValidRMC)),
					"Latitude":  Equal(NewFloat64(51.70021466666667)),
					"Longitude": Equal(NewFloat64(4.866866166666667)),
					"Speed":     Equal(NewFloat64(0)),
					"Course":    Equal(NewFloat64(345.6)),
					"Date":      Equal(NewDate(21, 4, 23)),
					"Variation": Equal(NewFloat64(0.3)),
				}))
			})
		})
		Context("a sentence with non existing status", func() {
			BeforeEach(func() {
				raw = "$GNRMC,220516,D,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*6B"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RMC struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":      Equal(NewTime(22, 5, 16, 0)),
					"Validity":  Equal(NewInvalidString("not a valid option")),
					"Latitude":  Equal(NewFloat64(51.56366666666666)),
					"Longitude": Equal(NewFloat64(-0.7040000000000001)),
					"Speed":     Equal(NewFloat64(173.8)),
					"Course":    Equal(NewFloat64(231.8)),
					"Date":      Equal(NewDate(94, 6, 13)),
					"Variation": Equal(NewFloat64(-4.2)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPRMC,114509.30,A,5142.01288,N,00452.01197,E,0.0,345.6,230421,0.3,E,A,C*12"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [5C != 12]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting directions from a $__RMC sentence", func() {
		BeforeEach(func() {
			parsed = RMC{
				Time: Time{
					Valid:       true,
					Hour:        20,
					Minute:      05,
					Second:      45,
					Millisecond: 315,
				},
				Validity:  NewString(ValidRMC),
				Latitude:  NewFloat64(Latitude),
				Longitude: NewFloat64(Longitude),
				Speed:     NewFloat64(SpeedOverGroundKnots),
				Course:    NewFloat64(TrueDirectionDegrees),
				Variation: NewFloat64(MagneticVariationDegrees),
				Date: Date{
					Valid: true,
					YY:    2021,
					MM:    4,
					DD:    16,
				},
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(Equal(Latitude))
				Expect(lon).To(Equal(Longitude))
			})
			It("returns a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("returns a valid magnetic variation", func() {
				Expect(parsed.GetMagneticVariation()).To(Float64Equal(MagneticVariationRadians, 0.00001))
			})
			It("returns a valid date and time", func() {
				Expect(parsed.GetDateTime()).To(Equal("2021-04-16T20:05:45.315Z"))
			})
		})
		Context("when having a struct with the validity flag set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Validity = NewString(InvalidRMC)
			})
			Specify("an error is returned when trying to retrieve the true course over ground", func() {
				value, err := parsed.GetTrueCourseOverGround()
				Expect(value).To(BeZero())
				Expect(err).To(HaveOccurred())
			})
			Specify("an error is returned when trying to retrieve the magnetic variation", func() {
				value, err := parsed.GetMagneticVariation()
				Expect(value).To(BeZero())
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
