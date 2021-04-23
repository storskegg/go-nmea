package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("GSV", func() {
	var (
		sentence Sentence
		parsed   GSV
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(GSV)
			} else {
				parsed = GSV{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*6B"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a valid sentence short sentence", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12*56"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
					}),
				}))
			})
		})
		Context("a sentence with an invalid numbers SVs in view", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11.2,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*77"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInvalidInt64("strconv.ParseInt: parsing \"11.2\": invalid syntax")),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a sentence with an invalid numbers number of messages", func() {
			BeforeEach(func() {
				raw = "$GLGSV,A3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInvalidInt64("strconv.ParseInt: parsing \"A3\": invalid syntax")),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a sentence with an invalid message number", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,A1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInvalidInt64("strconv.ParseInt: parsing \"A1\": invalid syntax")),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a sentence with an invalid SV prn number", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11,A03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInvalidInt64("strconv.ParseInt: parsing \"A03\": invalid syntax"), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a sentence with an invalid elevation", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11,03,A03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInvalidInt64("strconv.ParseInt: parsing \"A03\": invalid syntax"), Azimuth: NewInt64(111), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a sentence with an invalid azimuth", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11,03,03,A111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInvalidInt64("strconv.ParseInt: parsing \"A111\": invalid syntax"), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a sentence with an invalid SNR", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11,03,03,111,A00,04,15,270,00,06,01,010,12,13,06,292,00*2A"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GSV struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"TotalMessages":   Equal(NewInt64(3)),
					"MessageNumber":   Equal(NewInt64(1)),
					"NumberSVsInView": Equal(NewInt64(11)),
					"Info": Equal([]GSVInfo{
						{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInvalidInt64("strconv.ParseInt: parsing \"A00\": invalid syntax")},
						{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
						{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
						{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
					}),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GLGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*23"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [6B != 23]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a GSV struct", func() {
		BeforeEach(func() {
			parsed = GSV{
				NumberSVsInView: NewInt64(Satellites),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
		})
		Context("When having a parsed sentence without a number of satellites", func() {
			JustBeforeEach(func() {
				parsed.NumberSVsInView = NewInvalidInt64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetNumberOfSatellites()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
