package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("GGA", func() {
	var (
		sentence Sentence
		parsed   GGA
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(GGA)
			} else {
				parsed = GGA{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GNGGA,203415.000,6325.6138,N,01021.4290,E,1,8,2.42,72.5,M,41.5,M,,*7C"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GGA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":          Equal(NewTime(20, 34, 15, 0)),
					"Latitude":      Equal(NewFloat64(63.42689666666667)),
					"Longitude":     Equal(NewFloat64(10.357149999999999)),
					"FixQuality":    Equal(NewString("1")),
					"NumSatellites": Equal(NewInt64(8)),
					"HDOP":          Equal(NewFloat64(2.42)),
					"Altitude":      Equal(NewFloat64(72.5)),
					"Separation":    Equal(NewFloat64(41.5)),
					"DGPSAge":       Equal(NewString("")),
					"DGPSId":        Equal(NewString("")),
				}))
			})
		})
		Context("a sentence with a bad latitude", func() {
			BeforeEach(func() {
				raw = "$GNGGA,034225.077,A,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*24"
			})
			Specify("no error is returned", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GGA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":          Equal(NewTime(3, 42, 25, 77)),
					"Latitude":      Equal(NewInvalidFloat64("parse error (not decimal coordinate)")),
					"Longitude":     Equal(NewFloat64(151.40927833333333)),
					"FixQuality":    Equal(NewString("1")),
					"NumSatellites": Equal(NewInt64(3)),
					"HDOP":          Equal(NewFloat64(9.7)),
					"Altitude":      Equal(NewFloat64(-25)),
					"Separation":    Equal(NewFloat64(21)),
					"DGPSAge":       Equal(NewString("")),
					"DGPSId":        Equal(NewString("0000")),
				}))
			})
		})
		Context("a sentence with a bad longitude", func() {
			BeforeEach(func() {
				raw = "$GNGGA,034225.077,3356.4650,S,A,E,1,03,9.7,-25.0,M,21.0,M,,0000*12"
			})
			Specify("no error is returned", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GGA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":          Equal(NewTime(3, 42, 25, 77)),
					"Latitude":      Equal(NewFloat64(-33.94108333333334)),
					"Longitude":     Equal(NewInvalidFloat64("parse error (not decimal coordinate)")),
					"FixQuality":    Equal(NewString("1")),
					"NumSatellites": Equal(NewInt64(3)),
					"HDOP":          Equal(NewFloat64(9.7)),
					"Altitude":      Equal(NewFloat64(-25)),
					"Separation":    Equal(NewFloat64(21)),
					"DGPSAge":       Equal(NewString("")),
					"DGPSId":        Equal(NewString("0000")),
				}))
			})
		})
		Context("a sentence with a bad fix quality", func() {
			BeforeEach(func() {
				raw = "$GNGGA,034225.077,3356.4650,S,15124.5567,E,12,03,9.7,-25.0,M,21.0,M,,0000*7D"
			})
			Specify("no error is returned", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GGA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":          Equal(NewTime(3, 42, 25, 77)),
					"Latitude":      Equal(NewFloat64(-33.94108333333334)),
					"Longitude":     Equal(NewFloat64(151.40927833333333)),
					"FixQuality":    Equal(NewInvalidString("not a valid option")),
					"NumSatellites": Equal(NewInt64(3)),
					"HDOP":          Equal(NewFloat64(9.7)),
					"Altitude":      Equal(NewFloat64(-25)),
					"Separation":    Equal(NewFloat64(21)),
					"DGPSAge":       Equal(NewString("")),
					"DGPSId":        Equal(NewString("0000")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$SDGGA,0.5,0.5,*AA"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [7A != AA]"))
			})
		})
	})
	Describe("Getting data from a GGA struct", func() {
		BeforeEach(func() {
			parsed = GGA{
				Time:          Time{},
				Latitude:      NewFloat64(Latitude),
				Longitude:     NewFloat64(Longitude),
				FixQuality:    NewString(DGPS),
				NumSatellites: NewInt64(Satellites),
				HDOP:          NewInvalidFloat64(""),
				Altitude:      NewFloat64(Altitude),
				Separation:    NewInvalidFloat64(""),
				DGPSAge:       NewString(""),
				DGPSId:        NewString(""),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid position", func() {
				lat, lon, alt, _ := parsed.GetPosition3D()
				Expect(lat).To(Equal(Latitude))
				Expect(lon).To(Equal(Longitude))
				Expect(alt).To(Equal(Altitude))
			})
			It("returns a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("returns a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
		Context("when having a struct with a bad fix", func() {
			JustBeforeEach(func() {
				parsed.FixQuality = NewString(Invalid)
			})
			It("returns an error", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("returns a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(Invalid))
			})
		})
		Context("when having a struct with missing longitude", func() {
			JustBeforeEach(func() {
				parsed.Longitude = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("returns a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
		Context("when having a struct with missing latitude", func() {
			JustBeforeEach(func() {
				parsed.Latitude = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("returns a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
		Context("when having a struct with missing altitude", func() {
			JustBeforeEach(func() {
				parsed.Altitude = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("returns a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
		Context("when having a struct with missing number of satellites and fix quality", func() {
			JustBeforeEach(func() {
				parsed.NumSatellites = NewInvalidInt64("")
				parsed.FixQuality = NewInvalidString("")
			})
			It("returns an error", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetNumberOfSatellites()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetFixQuality()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
