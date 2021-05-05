package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("ZDA", func() {
	var (
		sentence Sentence
		parsed   ZDA
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(ZDA)
			} else {
				parsed = ZDA{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$GPZDA,172809.456,12,07,1996,00,00*57"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid ZDA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":          Equal(NewTime(17, 28, 9, 456)),
					"Day":           Equal(NewInt64(12)),
					"Month":         Equal(NewInt64(7)),
					"Year":          Equal(NewInt64(1996)),
					"OffsetHours":   Equal(NewInt64(0)),
					"OffsetMinutes": Equal(NewInt64(0)),
				}))
			})
		})
		Context("a sentence with an invalid date", func() {
			BeforeEach(func() {
				raw = "$GPZDA,172809.456,D,M,Y,00,00*04"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid ZDA struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":          Equal(NewTime(17, 28, 9, 456)),
					"Day":           Equal(NewInvalidInt64("strconv.ParseInt: parsing \"D\": invalid syntax")),
					"Month":         Equal(NewInvalidInt64("strconv.ParseInt: parsing \"M\": invalid syntax")),
					"Year":          Equal(NewInvalidInt64("strconv.ParseInt: parsing \"Y\": invalid syntax")),
					"OffsetHours":   Equal(NewInt64(0)),
					"OffsetMinutes": Equal(NewInt64(0)),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$GPZDA,172809.456,12,07,1996,00,00*BA"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [57 != BA]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a $__ZDA sentence", func() {
		BeforeEach(func() {
			parsed = ZDA{
				Time: Time{
					Valid:       true,
					Hour:        20,
					Minute:      05,
					Second:      45,
					Millisecond: 315,
				},
				Day:   NewInt64(16),
				Month: NewInt64(4),
				Year:  NewInt64(2021),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid date and time", func() {
				Expect(parsed.GetDateTime()).To(Equal("2021-04-16T20:05:45.315Z"))
			})
		})
		Context("when missing time", func() {
			JustBeforeEach(func() {
				parsed.Time = NewInvalidTime("")
			})
			It("returns an error", func() {
				_, err := parsed.GetDateTime()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when missing year", func() {
			JustBeforeEach(func() {
				parsed.Year = NewInvalidInt64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetDateTime()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when missing month", func() {
			JustBeforeEach(func() {
				parsed.Month = NewInvalidInt64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetDateTime()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when missing day", func() {
			JustBeforeEach(func() {
				parsed.Day = NewInvalidInt64("")
			})
			It("returns an error", func() {
				_, err := parsed.GetDateTime()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
