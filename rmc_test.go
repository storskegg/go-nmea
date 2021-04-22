package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var rmctests = []struct {
	name string
	raw  string
	err  string
	msg  RMC
}{
	{
		name: "good sentence A",
		raw:  "$GNRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*6E",
		msg: RMC{
			Time:      Time{true, 22, 05, 16, 0},
			Validity:  "A",
			Speed:     NewFloat64(173.8),
			Course:    NewFloat64(231.8),
			Date:      Date{true, 13, 06, 94},
			Variation: NewFloat64(-4.2),
			Latitude:  MustParseGPS("5133.82 N"),
			Longitude: MustParseGPS("00042.24 W"),
		},
	},
	{
		name: "good sentence B",
		raw:  "$GNRMC,142754.0,A,4302.539570,N,07920.379823,W,0.0,,070617,0.0,E,A*21",
		msg: RMC{
			Time:      Time{true, 14, 27, 54, 0},
			Validity:  "A",
			Speed:     NewFloat64(0),
			Course:    Float64{},
			Date:      Date{true, 7, 6, 17},
			Variation: NewFloat64(0),
			Latitude:  MustParseGPS("4302.539570 N"),
			Longitude: MustParseGPS("07920.379823 W"),
		},
	},
	{
		name: "good sentence C",
		raw:  "$GNRMC,100538.00,A,5546.27711,N,03736.91144,E,0.061,,260318,,,A*60",
		msg: RMC{
			Time:      Time{true, 10, 5, 38, 0},
			Validity:  "A",
			Speed:     NewFloat64(0.061),
			Course:    Float64{},
			Date:      Date{true, 26, 3, 18},
			Variation: Float64{},
			Latitude:  MustParseGPS("5546.27711 N"),
			Longitude: MustParseGPS("03736.91144 E"),
		},
	},
	{
		name: "bad sentence",
		raw:  "$GNRMC,220516,D,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*6B",
		err:  "nmea: GNRMC invalid validity: D",
	},
	{
		name: "good sentence A",
		raw:  "$GPRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*70",
		msg: RMC{
			Time:      Time{true, 22, 5, 16, 0},
			Validity:  "A",
			Speed:     NewFloat64(173.8),
			Course:    NewFloat64(231.8),
			Date:      Date{true, 13, 6, 94},
			Variation: NewFloat64(-4.2),
			Latitude:  MustParseGPS("5133.82 N"),
			Longitude: MustParseGPS("00042.24 W"),
		},
	},
	{
		name: "good sentence B",
		raw:  "$GPRMC,142754.0,A,4302.539570,N,07920.379823,W,0.0,,070617,0.0,E,A*3F",
		msg: RMC{
			Time:      Time{true, 14, 27, 54, 0},
			Validity:  "A",
			Speed:     NewFloat64(0),
			Course:    Float64{},
			Date:      Date{true, 7, 6, 17},
			Variation: NewFloat64(0),
			Latitude:  MustParseGPS("4302.539570 N"),
			Longitude: MustParseGPS("07920.379823 W"),
		},
	},
	{
		name: "bad validity",
		raw:  "$GPRMC,220516,D,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*75",
		err:  "nmea: GPRMC invalid validity: D",
	},
}

func TestRMC(t *testing.T) {
	for _, tt := range rmctests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				rmc := m.(RMC)
				rmc.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, rmc)
			}
		})
	}
}

var _ = Describe("RMC", func() {
	var (
		parsed RMC
	)
	Describe("Getting directions from a $__RMC sentence", func() {
		BeforeEach(func() {
			parsed = RMC{
				Time: Time{
					Valid:       true,
					Hour:        20,
					Minute:      05,
					Second:      45,
					Millisecond: 315,
				},
				Validity:  ValidRMC,
				Latitude:  NewFloat64(Latitude),
				Longitude: NewFloat64(Longitude),
				Speed:     NewFloat64(SpeedOverGroundKnots),
				Course:    NewFloat64(TrueDirectionDegrees),
				Variation: NewFloat64(MagneticVariationDegrees),
				Date: Date{
					Valid: true,
					YY:    2021,
					MM:    4,
					DD:    16,
				},
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(Equal(Latitude))
				Expect(lon).To(Equal(Longitude))
			})
			It("should give a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic variation", func() {
				Expect(parsed.GetMagneticVariation()).To(Float64Equal(MagneticVariationRadians, 0.00001))
			})
			It("should give a valid date and time", func() {
				Expect(parsed.GetDateTime()).To(Equal("2021-04-16T20:05:45.315Z"))
			})
		})
		Context("When having a parsed sentence with the validity flag set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Validity = InvalidRMC
			})
			Specify("an error is returned when trying to retrieve the true course over ground", func() {
				value, err := parsed.GetTrueCourseOverGround()
				Expect(value).To(BeZero())
				Expect(err).To(HaveOccurred())
			})
			Specify("an error is returned when trying to retrieve the magnetic variation", func() {
				value, err := parsed.GetMagneticVariation()
				Expect(value).To(BeZero())
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
