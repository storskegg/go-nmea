package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var ggatests = []struct {
	name string
	raw  string
	err  string
	msg  GGA
}{
	{
		name: "good sentence",
		raw:  "$GNGGA,203415.000,6325.6138,N,01021.4290,E,1,8,2.42,72.5,M,41.5,M,,*7C",
		msg: GGA{
			Time: Time{
				Valid:       true,
				Hour:        20,
				Minute:      34,
				Second:      15,
				Millisecond: 0,
			},
			Latitude:      MustParseLatLong("6325.6138 N"),
			Longitude:     MustParseLatLong("01021.4290 E"),
			FixQuality:    "1",
			NumSatellites: NewInt64(8),
			HDOP:          NewFloat64(2.42),
			Altitude:      NewFloat64(72.5),
			Separation:    NewFloat64(41.5),
			DGPSAge:       "",
			DGPSId:        "",
		},
	},
	{
		name: "bad latitude",
		raw:  "$GNGGA,034225.077,A,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*24",
		err:  "nmea: GNGGA invalid latitude: cannot parse [A S], unknown format",
	},
	{
		name: "bad longitude",
		raw:  "$GNGGA,034225.077,3356.4650,S,A,E,1,03,9.7,-25.0,M,21.0,M,,0000*12",
		err:  "nmea: GNGGA invalid longitude: cannot parse [A E], unknown format",
	},
	{
		name: "bad fix quality",
		raw:  "$GNGGA,034225.077,3356.4650,S,15124.5567,E,12,03,9.7,-25.0,M,21.0,M,,0000*7D",
		err:  "nmea: GNGGA invalid fix quality: 12",
	},
	{
		name: "good sentence",
		raw:  "$GPGGA,034225.077,3356.4650,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*51",
		msg: GGA{
			Time:          Time{true, 3, 42, 25, 77},
			Latitude:      MustParseLatLong("3356.4650 S"),
			Longitude:     MustParseLatLong("15124.5567 E"),
			FixQuality:    GPS,
			NumSatellites: NewInt64(03),
			HDOP:          NewFloat64(9.7),
			Altitude:      NewFloat64(-25.0),
			Separation:    NewFloat64(21.0),
			DGPSAge:       "",
			DGPSId:        "0000",
		},
	},
	{
		name: "bad latitude",
		raw:  "$GPGGA,034225.077,A,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*3A",
		err:  "nmea: GPGGA invalid latitude: cannot parse [A S], unknown format",
	},
	{
		name: "bad longitude",
		raw:  "$GPGGA,034225.077,3356.4650,S,A,E,1,03,9.7,-25.0,M,21.0,M,,0000*0C",
		err:  "nmea: GPGGA invalid longitude: cannot parse [A E], unknown format",
	},
	{
		name: "bad fix quality",
		raw:  "$GPGGA,034225.077,3356.4650,S,15124.5567,E,12,03,9.7,-25.0,M,21.0,M,,0000*63",
		err:  "nmea: GPGGA invalid fix quality: 12",
	},
}

func TestGGA(t *testing.T) {
	for _, tt := range ggatests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gga := m.(GGA)
				gga.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gga)
			}
		})
	}
}

var _ = Describe("GGA", func() {
	var (
		parsed GGA
	)
	Describe("Getting data from a $__GGA sentence", func() {
		BeforeEach(func() {
			parsed = GGA{
				Time:          Time{},
				Latitude:      NewFloat64(Latitude),
				Longitude:     NewFloat64(Longitude),
				FixQuality:    DGPS,
				NumSatellites: NewInt64(Satellites),
				HDOP:          Float64{},
				Altitude:      NewFloat64(Altitude),
				Separation:    Float64{},
				DGPSAge:       "",
				DGPSId:        "",
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid position", func() {
				lat, lon, alt, _ := parsed.GetPosition3D()
				Expect(lat).To(Equal(Latitude))
				Expect(lon).To(Equal(Longitude))
				Expect(alt).To(Equal(Altitude))
			})
			It("should give a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("should give a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
		Context("When having a parsed sentence with a bad fix", func() {
			JustBeforeEach(func() {
				parsed.FixQuality = Invalid
			})
			Specify("an error is returned", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("should give a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(Invalid))
			})
		})
		Context("When having a parsed sentence with missing longitude", func() {
			JustBeforeEach(func() {
				parsed.Longitude = Float64{}
			})
			Specify("an error is returned", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("should give a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
		Context("When having a parsed sentence with missing latitude", func() {
			JustBeforeEach(func() {
				parsed.Latitude = Float64{}
			})
			Specify("an error is returned", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("should give a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
		Context("When having a parsed sentence with missing altitude", func() {
			JustBeforeEach(func() {
				parsed.Altitude = Float64{}
			})
			Specify("an error is returned", func() {
				_, _, _, err := parsed.GetPosition3D()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
			It("should give a valid fix quality", func() {
				Expect(parsed.GetFixQuality()).To(Equal(DGPS))
			})
		})
	})
})
