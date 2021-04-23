// import (
// 	"testing"

// 	. "github.com/munnik/go-nmea"
// 	"github.com/stretchr/testify/assert"
// )

// var mtktests = []struct {
// 	name string
// 	raw  string
// 	err  string
// 	msg  MTK
// }{
// 	{
// 		name: "good: Packet Type: 001 PMTK_ACK",
// 		raw:  "$PMTK001,604,3*" + Checksum("PMTK001,604,3"),
// 		msg: MTK{
// 			Cmd:  NewInt64(604),
// 			Flag: NewInt64(3),
// 		},
// 	},
// 	{
// 		name: "missing flag",
// 		raw:  "$PMTK001,604*" + Checksum("PMTK001,604"),
// 		err:  "nmea: PMTK001 invalid flag: index out of range",
// 	},
// 	{
// 		name: "missing cmd",
// 		raw:  "$PMTK001*" + Checksum("PMTK001"),
// 		err:  "nmea: PMTK001 invalid command: index out of range",
// 	},
// }

// func TestMTK(t *testing.T) {
// 	for _, tt := range mtktests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			m, err := Parse(tt.raw)
// 			if tt.err != "" {
// 				assert.Error(t, err)
// 				assert.EqualError(t, err, tt.err)
// 			} else {
// 				assert.NoError(t, err)
// 				mtk := m.(MTK)
// 				mtk.BaseSentence = BaseSentence{}
// 				assert.Equal(t, tt.msg, mtk)
// 			}
// 		})
// 	}
// }

package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("MTK", func() {
	var (
		sentence Sentence
		parsed   MTK
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(MTK)
			} else {
				parsed = MTK{}
			}
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$PMTK001,604,3*32"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MTK struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Cmd":  Equal(NewInt64(604)),
					"Flag": Equal(NewInt64(3)),
				}))
			})
		})
		Context("a sentence missing flag", func() {
			BeforeEach(func() {
				raw = "$PMTK001,604*2D"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MTK struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Cmd":  Equal(NewInt64(604)),
					"Flag": Equal(NewInvalidInt64("index out of range")),
				}))
			})
		})
		Context("a sentence missing cmd and flag", func() {
			BeforeEach(func() {
				raw = "$PMTK001*33"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid MTK struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"Cmd":  Equal(NewInvalidInt64("index out of range")),
					"Flag": Equal(NewInvalidInt64("index out of range")),
				}))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$PMTK001,604,3*FF"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [32 != FF]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
})
