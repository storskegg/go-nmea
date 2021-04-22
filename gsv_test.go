package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var gsvtests = []struct {
	name string
	raw  string
	err  string
	msg  GSV
}{
	{
		name: "good sentence",
		raw:  "$GLGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*6B",
		msg: GSV{
			TotalMessages:   NewInt64(3),
			MessageNumber:   NewInt64(1),
			NumberSVsInView: NewInt64(11),
			Info: []GSVInfo{
				{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
				{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
			},
		},
	},
	{
		name: "short sentence",
		raw:  "$GLGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12*56",
		msg: GSV{
			TotalMessages:   NewInt64(3),
			MessageNumber:   NewInt64(1),
			NumberSVsInView: NewInt64(11),
			Info: []GSVInfo{
				{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
			},
		},
	},
	{
		name: "invalid number of svs",
		raw:  "$GLGSV,3,1,11.2,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*77",
		err:  "nmea: GLGSV invalid number of SVs in view: 11.2",
	},
	{
		name: "invalid number of messages",
		raw:  "$GLGSV,A3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A",
		err:  "nmea: GLGSV invalid total number of messages: A3",
	},
	{
		name: "invalid message number",
		raw:  "$GLGSV,3,A1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A",
		err:  "nmea: GLGSV invalid message number: A1",
	},
	{
		name: "invalid SV prn number",
		raw:  "$GLGSV,3,1,11,A03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A",
		err:  "nmea: GLGSV invalid SV prn number: A03",
	},
	{
		name: "invalid elevation",
		raw:  "$GLGSV,3,1,11,03,A03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A",
		err:  "nmea: GLGSV invalid elevation: A03",
	},
	{
		name: "invalid azimuth",
		raw:  "$GLGSV,3,1,11,03,03,A111,00,04,15,270,00,06,01,010,12,13,06,292,00*2A",
		err:  "nmea: GLGSV invalid azimuth: A111",
	},
	{
		name: "invalid SNR",
		raw:  "$GLGSV,3,1,11,03,03,111,A00,04,15,270,00,06,01,010,12,13,06,292,00*2A",
		err:  "nmea: GLGSV invalid SNR: A00",
	},
	{
		name: "good sentence",
		raw:  "$GPGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*77",
		msg: GSV{
			TotalMessages:   NewInt64(3),
			MessageNumber:   NewInt64(1),
			NumberSVsInView: NewInt64(11),
			Info: []GSVInfo{
				{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
				{SVPRNNumber: NewInt64(13), Elevation: NewInt64(6), Azimuth: NewInt64(292), SNR: NewInt64(0)},
			},
		},
	},
	{
		name: "short",
		raw:  "$GPGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12*4A",
		msg: GSV{
			TotalMessages:   NewInt64(3),
			MessageNumber:   NewInt64(1),
			NumberSVsInView: NewInt64(11),
			Info: []GSVInfo{
				{SVPRNNumber: NewInt64(3), Elevation: NewInt64(3), Azimuth: NewInt64(111), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(4), Elevation: NewInt64(15), Azimuth: NewInt64(270), SNR: NewInt64(0)},
				{SVPRNNumber: NewInt64(6), Elevation: NewInt64(1), Azimuth: NewInt64(10), SNR: NewInt64(12)},
			},
		},
	},
	{
		name: "invalid number of SVs",
		raw:  "$GPGSV,3,1,11.2,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*6b",
		err:  "nmea: GPGSV invalid number of SVs in view: 11.2",
	},
	{
		name: "invalid total number of messages",
		raw:  "$GPGSV,A3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*36",
		err:  "nmea: GPGSV invalid total number of messages: A3",
	},
	{
		name: "invalid message number",
		raw:  "$GPGSV,3,A1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*36",
		err:  "nmea: GPGSV invalid message number: A1",
	},
	{
		name: "invalid SV prn number",
		raw:  "$GPGSV,3,1,11,A03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*36",
		err:  "nmea: GPGSV invalid SV prn number: A03",
	},
	{
		name: "invalid elevation",
		raw:  "$GPGSV,3,1,11,03,A03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*36",
		err:  "nmea: GPGSV invalid elevation: A03",
	},
	{
		name: "invalid azimuth",
		raw:  "$GPGSV,3,1,11,03,03,A111,00,04,15,270,00,06,01,010,12,13,06,292,00*36",
		err:  "nmea: GPGSV invalid azimuth: A111",
	},
	{
		name: "invalid SNR",
		raw:  "$GPGSV,3,1,11,03,03,111,A00,04,15,270,00,06,01,010,12,13,06,292,00*36",
		err:  "nmea: GPGSV invalid SNR: A00",
	},
}

func TestGSV(t *testing.T) {
	for _, tt := range gsvtests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gsv := m.(GSV)
				gsv.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gsv)
			}
		})
	}
}

var _ = Describe("GSV", func() {
	var (
		parsed GSV
	)
	Describe("Getting data from a $__GSV sentence", func() {
		BeforeEach(func() {
			parsed = GSV{
				NumberSVsInView: NewInt64(Satellites),
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid number of satellites", func() {
				Expect(parsed.GetNumberOfSatellites()).To(Equal(Satellites))
			})
		})
		Context("When having a parsed sentence without a number of satellites", func() {
			JustBeforeEach(func() {
				parsed.NumberSVsInView = Int64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetNumberOfSatellites()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
