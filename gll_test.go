package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("GLL", func() {
	var (
		sentence Sentence
		parsed   GLL
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(GLL)
			} else {
				parsed = GLL{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPGLL,3926.7952,N,12000.5947,W,022732,A,A*58"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GLL struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":      Equal(NewTime(2, 27, 32, 0)),
					"Latitude":  Equal(NewFloat64(39.44658666666667)),
					"Longitude": Equal(NewFloat64(-120.00991166666667)),
					"Validity":  Equal(NewString(ValidGLL)),
				}))
			})
		})
		Context("a valid sentence with an invalid validity", func() {
			BeforeEach(func() {
				raw = "$GPGLL,3926.7952,N,12000.5947,W,022732,V,A*4F"
			})
			Specify("no error is returned", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GLL struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":      Equal(NewTime(2, 27, 32, 0)),
					"Latitude":  Equal(NewFloat64(39.44658666666667)),
					"Longitude": Equal(NewFloat64(-120.00991166666667)),
					"Validity":  Equal(NewString(InvalidGLL)),
				}))
			})
		})
		Context("a sentence with a non existing validity", func() {
			BeforeEach(func() {
				raw = "$GPGLL,3926.7952,N,12000.5947,W,022732,D,A*5D"
			})
			Specify("no error is returned", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid GLL struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":      Equal(NewTime(2, 27, 32, 0)),
					"Latitude":  Equal(NewFloat64(39.44658666666667)),
					"Longitude": Equal(NewFloat64(-120.00991166666667)),
					"Validity":  Equal(NewInvalidString("not a valid option")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPGLL,3926.7952,N,12000.5947,W,022732,A,A*BC"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [58 != BC]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a GLL struct", func() {
		BeforeEach(func() {
			parsed = GLL{
				Time:      Time{},
				Latitude:  NewFloat64(Latitude),
				Longitude: NewFloat64(Longitude),
				Validity:  NewString(ValidGLL),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(Equal(Latitude))
				Expect(lon).To(Equal(Longitude))
			})
		})
		Context("when having a struct with validity set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Validity = NewString(InvalidGLL)
			})
			It("returns an error", func() {
				_, _, err := parsed.GetPosition2D()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a struct with missing longitude", func() {
			JustBeforeEach(func() {
				parsed.Longitude = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, _, err := parsed.GetPosition2D()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a struct with missing latitude", func() {
			JustBeforeEach(func() {
				parsed.Latitude = NewInvalidFloat64("")
			})
			It("returns an error", func() {
				_, _, err := parsed.GetPosition2D()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
