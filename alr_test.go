package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("ALR", func() {
	var (
		sentence Sentence
		parsed   ALR
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(ALR)
			} else {
				parsed = ALR{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$AIALR,100615.00,002,V,V,AIS: Antenna VSWR exceeds limit*46"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid ALR struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Time":        Equal(NewTime(10, 6, 15, 0)),
					"Identifier":  Equal(NewString("002")),
					"Condition":   Equal(NewString(InactiveALR)),
					"State":       Equal(NewString(UnacknowledgedALR)),
					"Description": Equal(NewString("AIS: Antenna VSWR exceeds limit")),
				}))
			})
		})
	})
	Describe("Getting data from a $__ALR sentence", func() {
		BeforeEach(func() {
			parsed = ALR{
				Time: Time{
					Valid:       true,
					Hour:        20,
					Minute:      05,
					Second:      45,
					Millisecond: 315,
				},
				Identifier:  NewString("002"),
				Condition:   NewString(ActiveALR),
				State:       NewString(AcknowledgedALR),
				Description: NewString("Testing"),
			}
		})
		Context("when having a complete struct", func() {
			It("returns a valid identifier", func() {
				identifier, _ := parsed.GetIdentifier()
				Expect(identifier).To(Equal("002"))
			})
			It("returns a valid condition", func() {
				condition, _ := parsed.IsActive()
				Expect(condition).To(Equal(true))
			})
			It("returns a valid state", func() {
				state, _ := parsed.IsUnacknowledged()
				Expect(state).To(Equal(false))
			})
			It("returns a valid description", func() {
				description, _ := parsed.GetDescription()
				Expect(description).To(Equal("Testing"))
			})
		})
	})
})
