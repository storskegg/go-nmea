package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("WPL", func() {
	var (
		sentence Sentence
		parsed   WPL
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(WPL)
			} else {
				parsed = WPL{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$IIWPL,5503.4530,N,01037.2742,E,411*6F"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid WPL struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Latitude":  Equal(NewFloat64(55.057550000000006)),
					"Longitude": Equal(NewFloat64(10.621236666666668)),
					"Ident":     Equal(NewString("411")),
				}))
			})
		})
		Context("a valid sentence with a bad latitude", func() {
			BeforeEach(func() {
				raw = "$IIWPL,A,N,01037.2742,E,411*01"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid WPL struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Latitude":  Equal(NewInvalidFloat64("parse error (not decimal coordinate)")),
					"Longitude": Equal(NewFloat64(10.621236666666668)),
					"Ident":     Equal(NewString("411")),
				}))
			})
		})
		Context("a valid sentence with a bad longitude", func() {
			BeforeEach(func() {
				raw = "$IIWPL,5503.4530,N,A,E,411*36"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid WPL struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Latitude":  Equal(NewFloat64(55.057550000000006)),
					"Longitude": Equal(NewInvalidFloat64("parse error (not decimal coordinate)")),
					"Ident":     Equal(NewString("411")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$IIWPL,5503.4530,N,01037.2742,E,411*1C"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [6F != 1C]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
})
