package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("RTE", func() {
	var (
		sentence Sentence
		parsed   RTE
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(RTE)
			} else {
				parsed = RTE{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$IIRTE,4,1,c,Rte 1,411,412,413,414,415*6F"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RTE struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumberOfSentences":         Equal(NewInt64(4)),
					"SentenceNumber":            Equal(NewInt64(1)),
					"ActiveRouteOrWaypointList": Equal(NewString(ActiveRoute)),
					"Name":                      Equal(NewString("Rte 1")),
					"Idents":                    Equal(NewStringList([]String{NewString("411"), NewString("412"), NewString("413"), NewString("414"), NewString("415")})),
				}))
			})
		})
		Context("a sentence without any waypoints", func() {
			BeforeEach(func() {
				raw = "$IIRTE,4,1,c,Rte 1*77"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RTE struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumberOfSentences":         Equal(NewInt64(4)),
					"SentenceNumber":            Equal(NewInt64(1)),
					"ActiveRouteOrWaypointList": Equal(NewString(ActiveRoute)),
					"Name":                      Equal(NewString("Rte 1")),
					"Idents":                    Equal(NewInvalidStringList("index out of range")),
				}))
			})
		})
		Context("a sentence with an invalid number of sentences", func() {
			BeforeEach(func() {
				raw = "$IIRTE,X,1,c,Rte 1,411,412,413,414,415*03"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RTE struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumberOfSentences":         Equal(NewInvalidInt64("strconv.ParseInt: parsing \"X\": invalid syntax")),
					"SentenceNumber":            Equal(NewInt64(1)),
					"ActiveRouteOrWaypointList": Equal(NewString(ActiveRoute)),
					"Name":                      Equal(NewString("Rte 1")),
					"Idents":                    Equal(NewStringList([]String{NewString("411"), NewString("412"), NewString("413"), NewString("414"), NewString("415")})),
				}))
			})
		})
		Context("a sentence with an invalid sentence number", func() {
			BeforeEach(func() {
				raw = "$IIRTE,4,X,c,Rte 1,411,412,413,414,415*06"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RTE struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumberOfSentences":         Equal(NewInt64(4)),
					"SentenceNumber":            Equal(NewInvalidInt64("strconv.ParseInt: parsing \"X\": invalid syntax")),
					"ActiveRouteOrWaypointList": Equal(NewString(ActiveRoute)),
					"Name":                      Equal(NewString("Rte 1")),
					"Idents":                    Equal(NewStringList([]String{NewString("411"), NewString("412"), NewString("413"), NewString("414"), NewString("415")})),
				}))
			})
		})
		Context("a sentence with an invalid active route or waypoint list", func() {
			BeforeEach(func() {
				raw = "$IIRTE,4,1,X,Rte 1,411,412,413,414,415*54"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid RTE struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumberOfSentences":         Equal(NewInt64(4)),
					"SentenceNumber":            Equal(NewInt64(1)),
					"ActiveRouteOrWaypointList": Equal(NewInvalidString("not a valid option")),
					"Name":                      Equal(NewString("Rte 1")),
					"Idents":                    Equal(NewStringList([]String{NewString("411"), NewString("412"), NewString("413"), NewString("414"), NewString("415")})),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$IIRTE,4,1,c,Rte 1,411,412,413,414,415*70"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [6F != 70]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
})
