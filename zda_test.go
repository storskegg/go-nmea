package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var zdatests = []struct {
	name string
	raw  string
	err  string
	msg  ZDA
}{
	{
		name: "good sentence",
		raw:  "$GPZDA,172809.456,12,07,1996,00,00*57",
		msg: ZDA{
			Time: Time{
				Valid:       true,
				Hour:        17,
				Minute:      28,
				Second:      9,
				Millisecond: 456,
			},
			Day:           NewInt64(12),
			Month:         NewInt64(7),
			Year:          NewInt64(1996),
			OffsetHours:   NewInt64(0),
			OffsetMinutes: NewInt64(0),
		},
	},
	{
		name: "invalid day",
		raw:  "$GPZDA,220516,D,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*76",
		err:  "nmea: GPZDA invalid day: D",
	},
}

func TestZDA(t *testing.T) {
	for _, tt := range zdatests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				zda := m.(ZDA)
				zda.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, zda)
			}
		})
	}
}

var _ = Describe("ZDA", func() {
	var (
		parsed ZDA
	)
	Describe("Parse an actual sentence", func() {
		BeforeEach(func() {
			sentence := "$GPZDA,185257,17,04,2021,00,00*47"
			parseResult, err := Parse(sentence)
			if err != nil {
				Fail("Could not parse sentence")
			}
			var ok bool
			if parsed, ok = parseResult.(ZDA); !ok {
				Fail("Could not cast to ZDA")
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid date and time", func() {
				Expect(parsed.GetDateTime()).To(Equal("2021-04-17T18:52:57Z"))
			})
		})
	})
	Describe("Getting directions from a $__ZDA sentence", func() {
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
		Context("When having a parsed sentence", func() {
			It("should give a valid date and time", func() {
				Expect(parsed.GetDateTime()).To(Equal("2021-04-16T20:05:45.315Z"))
			})
		})
	})
})
