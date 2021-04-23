package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("GNS", func() {
	var (
		sentence Sentence
		parsed   GNS
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(GNS)
			} else {
				parsed = GNS{}
			}
		})
		Context("a valid sentence with mode RR", func() {
			BeforeEach(func() {
				raw = "$GNGNS,014035.00,4332.69262,S,17235.48549,E,RR,13,0.9,25.63,11.24,,*70"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GNS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":       Equal(NewTime(1, 40, 35, 0)),
					"Latitude":   Equal(NewFloat64(-43.544877)),
					"Longitude":  Equal(NewFloat64(172.59142483333332)),
					"Mode":       Equal([]String{NewString("R"), NewString("R")}),
					"SVs":        Equal(NewInt64(13)),
					"HDOP":       Equal(NewFloat64(0.9)),
					"Altitude":   Equal(NewFloat64(25.63)),
					"Separation": Equal(NewFloat64(11.24)),
					"Age":        Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"Station":    Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
				}))
			})
		})
		Context("a valid sentence with mode AA", func() {
			BeforeEach(func() {
				raw = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,AA,14,0.6,161.5,48.0,,*6D"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GNS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":       Equal(NewTime(9, 48, 21, 0)),
					"Latitude":   Equal(NewFloat64(48.83218845)),
					"Longitude":  Equal(NewFloat64(2.2675553833333333)),
					"Mode":       Equal([]String{NewString("A"), NewString("A")}),
					"SVs":        Equal(NewInt64(14)),
					"HDOP":       Equal(NewFloat64(0.6)),
					"Altitude":   Equal(NewFloat64(161.5)),
					"Separation": Equal(NewFloat64(48)),
					"Age":        Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"Station":    Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
				}))
			})
		})
		Context("a valid sentence with mode AAN", func() {
			BeforeEach(func() {
				raw = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,AAN,14,0.6,161.5,48.0,,*23"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GNS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":       Equal(NewTime(9, 48, 21, 0)),
					"Latitude":   Equal(NewFloat64(48.83218845)),
					"Longitude":  Equal(NewFloat64(2.2675553833333333)),
					"Mode":       Equal([]String{NewString("A"), NewString("A"), NewString("N")}),
					"SVs":        Equal(NewInt64(14)),
					"HDOP":       Equal(NewFloat64(0.6)),
					"Altitude":   Equal(NewFloat64(161.5)),
					"Separation": Equal(NewFloat64(48)),
					"Age":        Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"Station":    Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
				}))
			})
		})
		Context("a valid sentence with an empty mode", func() {
			BeforeEach(func() {
				raw = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,,14,0.6,161.5,48.0,,*6D"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GNS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":       Equal(NewTime(9, 48, 21, 0)),
					"Latitude":   Equal(NewFloat64(48.83218845)),
					"Longitude":  Equal(NewFloat64(2.2675553833333333)),
					"Mode":       Equal([]String{}),
					"SVs":        Equal(NewInt64(14)),
					"HDOP":       Equal(NewFloat64(0.6)),
					"Altitude":   Equal(NewFloat64(161.5)),
					"Separation": Equal(NewFloat64(48)),
					"Age":        Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"Station":    Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
				}))
			})
		})
		Context("a sentence with a non existing mode", func() {
			BeforeEach(func() {
				raw = "$GNGNS,094821.0,4849.931307,N,00216.053323,E,PXKR,14,0.6,161.5,48.0,,*7C"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GNS struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":       Equal(NewTime(9, 48, 21, 0)),
					"Latitude":   Equal(NewFloat64(48.83218845)),
					"Longitude":  Equal(NewFloat64(2.2675553833333333)),
					"Mode":       Equal([]String{}),
					"SVs":        Equal(NewInt64(14)),
					"HDOP":       Equal(NewFloat64(0.6)),
					"Altitude":   Equal(NewFloat64(161.5)),
					"Separation": Equal(NewFloat64(48)),
					"Age":        Equal(NewInvalidFloat64("strconv.ParseFloat: parsing \"\": invalid syntax")),
					"Station":    Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GNGNS,014035.00,4332.69262,S,17235.48549,E,RR,13,0.9,25.63,11.24,,*71"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [70 != 71]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a GNS struct", func() {
		BeforeEach(func() {
			parsed = GNS{
				Time:       Time{},
				Latitude:   NewFloat64(Latitude),
				Longitude:  NewFloat64(Longitude),
				Mode:       []String{NewString(SimulatorGNS)},
				SVs:        NewInvalidInt64(""),
				HDOP:       NewInvalidFloat64(""),
				Altitude:   NewFloat64(Altitude),
				Separation: NewInvalidFloat64(""),
				Age:        NewInvalidFloat64(""),
				Station:    NewInvalidInt64(""),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid position", func() {
				lat, lon, alt, _ := parsed.GetPosition3D()
				Expect(lat).To(Equal(Latitude))
				Expect(lon).To(Equal(Longitude))
				Expect(alt).To(Equal(Altitude))
			})
		})
	})
})
