package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var gsatests = []struct {
	name string
	raw  string
	err  string
	msg  GSA
}{
	{
		name: "good sentence",
		raw:  "$GPGSA,A,3,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*36",
		msg: GSA{
			Mode:    "A",
			FixType: "3",
			SV:      []string{"22", "19", "18", "27", "14", "03"},
			PDOP:    NewFloat64(3.1),
			HDOP:    NewFloat64(2),
			VDOP:    NewFloat64(2.4),
		},
	},
	{
		name: "bad mode",
		raw:  "$GPGSA,F,3,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*31",
		err:  "nmea: GPGSA invalid selection mode: F",
	},
	{
		name: "bad fix",
		raw:  "$GPGSA,A,6,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*33",
		err:  "nmea: GPGSA invalid fix type: 6",
	},
}

func TestGSA(t *testing.T) {
	for _, tt := range gsatests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gsa := m.(GSA)
				gsa.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gsa)
			}
		})
	}
}

var _ = Describe("GSA", func() {
	var (
		parsed GSA
	)
	Describe("Getting data from a $__GSA sentence", func() {
		BeforeEach(func() {
			parsed = GSA{
				Mode:    Auto,
				FixType: Fix3D,
				SV:      make([]string, Satellites),
				PDOP:    Float64{},
				HDOP:    Float64{},
				VDOP:    Float64{},
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("should give a valid fix type", func() {
				Expect(parsed.GetFixType()).To(Equal(Fix3D))
			})
		})
	})
})
