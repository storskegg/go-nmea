package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var hdttests = []struct {
	name string
	raw  string
	err  string
	msg  HDT
}{
	{
		name: "good sentence",
		raw:  "$GPHDT,123.456,T*32",
		msg: HDT{
			Heading: NewFloat64(123.456),
			True:    true,
		},
	},
	{
		name: "invalid True",
		raw:  "$GPHDT,123.456,X*3E",
		err:  "nmea: GPHDT invalid true: X",
	},
	{
		name: "invalid Heading",
		raw:  "$GPHDT,XXX,T*43",
		err:  "nmea: GPHDT invalid heading: XXX",
	},
}

func TestHDT(t *testing.T) {
	for _, tt := range hdttests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				hdt := m.(HDT)
				hdt.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, hdt)
			}
		})
	}
}

var _ = Describe("HDT", func() {
	var (
		parsed HDT
	)
	Describe("Getting data from a $__HDT sentence", func() {
		BeforeEach(func() {
			parsed = HDT{
				Heading: NewFloat64(TrueDirectionDegrees),
				True:    true,
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing heading", func() {
			JustBeforeEach(func() {
				parsed.Heading = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When having a parsed sentence with true flag not set", func() {
			JustBeforeEach(func() {
				parsed.True = false
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
