package nmea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var glgsvtests = []struct {
	name string
	raw  string
	err  string
	msg  GLGSV
}{
	{
		name: "good sentence",
		raw:  "$GLGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*6B",
		msg: GLGSV{
			TotalMessages:   Int64{Valid: true, Value: 3},
			MessageNumber:   Int64{Valid: true, Value: 1},
			NumberSVsInView: Int64{Valid: true, Value: 11},
			Info: []GLGSVInfo{
				{SVPRNNumber: Int64{Valid: true, Value: 3}, Elevation: Int64{Valid: true, Value: 3}, Azimuth: Int64{Valid: true, Value: 111}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 4}, Elevation: Int64{Valid: true, Value: 15}, Azimuth: Int64{Valid: true, Value: 270}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 6}, Elevation: Int64{Valid: true, Value: 1}, Azimuth: Int64{Valid: true, Value: 10}, SNR: Int64{Valid: true, Value: 12}},
				{SVPRNNumber: Int64{Valid: true, Value: 13}, Elevation: Int64{Valid: true, Value: 6}, Azimuth: Int64{Valid: true, Value: 292}, SNR: Int64{Valid: true, Value: 0}},
			},
		},
	},
	{
		name: "short sentence",
		raw:  "$GLGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12*56",
		msg: GLGSV{
			TotalMessages:   Int64{Valid: true, Value: 3},
			MessageNumber:   Int64{Valid: true, Value: 1},
			NumberSVsInView: Int64{Valid: true, Value: 11},
			Info: []GLGSVInfo{
				{SVPRNNumber: Int64{Valid: true, Value: 3}, Elevation: Int64{Valid: true, Value: 3}, Azimuth: Int64{Valid: true, Value: 111}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 4}, Elevation: Int64{Valid: true, Value: 15}, Azimuth: Int64{Valid: true, Value: 270}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 6}, Elevation: Int64{Valid: true, Value: 1}, Azimuth: Int64{Valid: true, Value: 10}, SNR: Int64{Valid: true, Value: 12}},
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
}

func TestGLGSV(t *testing.T) {
	for _, tt := range glgsvtests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				glgsv := m.(GLGSV)
				glgsv.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, glgsv)
			}
		})
	}
}

var gnggatests = []struct {
	name string
	raw  string
	err  string
	msg  GNGGA
}{
	{
		name: "good sentence",
		raw:  "$GNGGA,203415.000,6325.6138,N,01021.4290,E,1,8,2.42,72.5,M,41.5,M,,*7C",
		msg: GNGGA{
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
			NumSatellites: Int64{Valid: true, Value: 8},
			HDOP:          Float64{Valid: true, Value: 2.42},
			Altitude:      Float64{Valid: true, Value: 72.5},
			Separation:    Float64{Valid: true, Value: 41.5},
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
}

func TestGNGGA(t *testing.T) {
	for _, tt := range gnggatests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gngga := m.(GNGGA)
				gngga.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gngga)
			}
		})
	}
}

var gngnstests = []struct {
	name string
	raw  string
	err  string
	msg  GNGNS
}{
	{
		name: "good sentence A",
		raw:  "$GNGNS,014035.00,4332.69262,S,17235.48549,E,RR,13,0.9,25.63,11.24,,*70",
		msg: GNGNS{
			Time:       Time{true, 1, 40, 35, 0},
			Latitude:   MustParseGPS("4332.69262 S"),
			Longitude:  MustParseGPS("17235.48549 E"),
			Mode:       []string{"R", "R"},
			SVs:        Int64{Valid: true, Value: 13},
			HDOP:       Float64{Valid: true, Value: 0.9},
			Altitude:   Float64{Valid: true, Value: 25.63},
			Separation: Float64{Valid: true, Value: 11.24},
			Age:        Float64{Valid: false, Value: 0},
			Station:    Int64{Valid: false, Value: 0},
		},
	},
	{
		name: "good sentence B",
		raw:  "$GNGNS,094821.0,4849.931307,N,00216.053323,E,AA,14,0.6,161.5,48.0,,*6D",
		msg: GNGNS{
			Time:       Time{true, 9, 48, 21, 0},
			Latitude:   MustParseGPS("4849.931307 N"),
			Longitude:  MustParseGPS("00216.053323 E"),
			Mode:       []string{"A", "A"},
			SVs:        Int64{Valid: true, Value: 14},
			HDOP:       Float64{Valid: true, Value: 0.6},
			Altitude:   Float64{Valid: true, Value: 161.5},
			Separation: Float64{Valid: true, Value: 48.0},
			Age:        Float64{Valid: false, Value: 0},
			Station:    Int64{Valid: false, Value: 0},
		},
	},
	{
		name: "good sentence B",
		raw:  "$GNGNS,094821.0,4849.931307,N,00216.053323,E,AAN,14,0.6,161.5,48.0,,*23",
		msg: GNGNS{
			Time:       Time{true, 9, 48, 21, 0},
			Latitude:   MustParseGPS("4849.931307 N"),
			Longitude:  MustParseGPS("00216.053323 E"),
			Mode:       []string{"A", "A", "N"},
			SVs:        Int64{Valid: true, Value: 14},
			HDOP:       Float64{Valid: true, Value: 0.6},
			Altitude:   Float64{Valid: true, Value: 161.5},
			Separation: Float64{Valid: true, Value: 48.0},
			Age:        Float64{Valid: false, Value: 0},
			Station:    Int64{Valid: false, Value: 0},
		},
	},
	{
		name: "bad sentence",
		raw:  "$GNGNS,094821.0,4849.931307,N,00216.053323,E,AAX,14,0.6,161.5,48.0,,*35",
		err:  "nmea: GNGNS invalid mode: AAX",
	},
}

func TestGNGNS(t *testing.T) {
	for _, tt := range gngnstests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gngns := m.(GNGNS)
				gngns.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gngns)
			}
		})
	}
}

var gnrmctests = []struct {
	name string
	raw  string
	err  string
	msg  GNRMC
}{
	{
		name: "good sentence A",
		raw:  "$GNRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*6E",
		msg: GNRMC{
			Time:      Time{true, 22, 05, 16, 0},
			Validity:  "A",
			Speed:     Float64{Valid: true, Value: 173.8},
			Course:    Float64{Valid: true, Value: 231.8},
			Date:      Date{true, 13, 06, 94},
			Variation: Float64{Valid: true, Value: -4.2},
			Latitude:  MustParseGPS("5133.82 N"),
			Longitude: MustParseGPS("00042.24 W"),
		},
	},
	{
		name: "good sentence B",
		raw:  "$GNRMC,142754.0,A,4302.539570,N,07920.379823,W,0.0,,070617,0.0,E,A*21",
		msg: GNRMC{
			Time:      Time{true, 14, 27, 54, 0},
			Validity:  "A",
			Speed:     Float64{Valid: true, Value: 0},
			Course:    Float64{Valid: false, Value: 0},
			Date:      Date{true, 7, 6, 17},
			Variation: Float64{Valid: true, Value: 0},
			Latitude:  MustParseGPS("4302.539570 N"),
			Longitude: MustParseGPS("07920.379823 W"),
		},
	},
	{
		name: "good sentence C",
		raw:  "$GNRMC,100538.00,A,5546.27711,N,03736.91144,E,0.061,,260318,,,A*60",
		msg: GNRMC{
			Time:      Time{true, 10, 5, 38, 0},
			Validity:  "A",
			Speed:     Float64{Valid: true, Value: 0.061},
			Course:    Float64{Valid: false, Value: 0},
			Date:      Date{true, 26, 3, 18},
			Variation: Float64{Valid: false, Value: 0},
			Latitude:  MustParseGPS("5546.27711 N"),
			Longitude: MustParseGPS("03736.91144 E"),
		},
	},
	{
		name: "bad sentence",
		raw:  "$GNRMC,220516,D,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*6B",
		err:  "nmea: GNRMC invalid validity: D",
	},
}

func TestGNRMC(t *testing.T) {
	for _, tt := range gnrmctests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gnrmc := m.(GNRMC)
				gnrmc.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gnrmc)
			}
		})
	}
}

var gpggatests = []struct {
	name string
	raw  string
	err  string
	msg  GPGGA
}{
	{
		name: "good sentence",
		raw:  "$GPGGA,034225.077,3356.4650,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*51",
		msg: GPGGA{
			Time:          Time{true, 3, 42, 25, 77},
			Latitude:      MustParseLatLong("3356.4650 S"),
			Longitude:     MustParseLatLong("15124.5567 E"),
			FixQuality:    GPS,
			NumSatellites: Int64{Valid: true, Value: 03},
			HDOP:          Float64{Valid: true, Value: 9.7},
			Altitude:      Float64{Valid: true, Value: -25.0},
			Separation:    Float64{Valid: true, Value: 21.0},
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

func TestGPGGA(t *testing.T) {
	for _, tt := range gpggatests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gpgga := m.(GPGGA)
				gpgga.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gpgga)
			}
		})
	}
}

var gpglltests = []struct {
	name string
	raw  string
	err  string
	msg  GPGLL
}{
	{
		name: "good sentence",
		raw:  "$GPGLL,3926.7952,N,12000.5947,W,022732,A,A*58",
		msg: GPGLL{
			Latitude:  MustParseLatLong("3926.7952 N"),
			Longitude: MustParseLatLong("12000.5947 W"),
			Time: Time{
				Valid:       true,
				Hour:        2,
				Minute:      27,
				Second:      32,
				Millisecond: 0,
			},
			Validity: "A",
		},
	},
	{
		name: "bad validity",
		raw:  "$GPGLL,3926.7952,N,12000.5947,W,022732,D,A*5D",
		err:  "nmea: GPGLL invalid validity: D",
	},
}

func TestGPGLL(t *testing.T) {
	for _, tt := range gpglltests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gpgll := m.(GPGLL)
				gpgll.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gpgll)
			}
		})
	}
}

var gpgsatests = []struct {
	name string
	raw  string
	err  string
	msg  GPGSA
}{
	{
		name: "good sentence",
		raw:  "$GPGSA,A,3,22,19,18,27,14,03,,,,,,,3.1,2.0,2.4*36",
		msg: GPGSA{
			Mode:    "A",
			FixType: "3",
			SV:      []string{"22", "19", "18", "27", "14", "03"},
			PDOP:    Float64{Valid: true, Value: 3.1},
			HDOP:    Float64{Valid: true, Value: 2},
			VDOP:    Float64{Valid: true, Value: 2.4},
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

func TestGPGSA(t *testing.T) {
	for _, tt := range gpgsatests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gpgsa := m.(GPGSA)
				gpgsa.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gpgsa)
			}
		})
	}
}

var gpgsvtests = []struct {
	name string
	raw  string
	err  string
	msg  GPGSV
}{
	{
		name: "good sentence",
		raw:  "$GPGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12,13,06,292,00*77",
		msg: GPGSV{
			TotalMessages:   Int64{Valid: true, Value: 3},
			MessageNumber:   Int64{Valid: true, Value: 1},
			NumberSVsInView: Int64{Valid: true, Value: 11},
			Info: []GPGSVInfo{
				{SVPRNNumber: Int64{Valid: true, Value: 3}, Elevation: Int64{Valid: true, Value: 3}, Azimuth: Int64{Valid: true, Value: 111}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 4}, Elevation: Int64{Valid: true, Value: 15}, Azimuth: Int64{Valid: true, Value: 270}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 6}, Elevation: Int64{Valid: true, Value: 1}, Azimuth: Int64{Valid: true, Value: 10}, SNR: Int64{Valid: true, Value: 12}},
				{SVPRNNumber: Int64{Valid: true, Value: 13}, Elevation: Int64{Valid: true, Value: 6}, Azimuth: Int64{Valid: true, Value: 292}, SNR: Int64{Valid: true, Value: 0}},
			},
		},
	},
	{
		name: "short",
		raw:  "$GPGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,12*4A",
		msg: GPGSV{
			TotalMessages:   Int64{Valid: true, Value: 3},
			MessageNumber:   Int64{Valid: true, Value: 1},
			NumberSVsInView: Int64{Valid: true, Value: 11},
			Info: []GPGSVInfo{
				{SVPRNNumber: Int64{Valid: true, Value: 3}, Elevation: Int64{Valid: true, Value: 3}, Azimuth: Int64{Valid: true, Value: 111}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 4}, Elevation: Int64{Valid: true, Value: 15}, Azimuth: Int64{Valid: true, Value: 270}, SNR: Int64{Valid: true, Value: 0}},
				{SVPRNNumber: Int64{Valid: true, Value: 6}, Elevation: Int64{Valid: true, Value: 1}, Azimuth: Int64{Valid: true, Value: 10}, SNR: Int64{Valid: true, Value: 12}},
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

func TestGPGSV(t *testing.T) {
	for _, tt := range gpgsvtests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gpgsv := m.(GPGSV)
				gpgsv.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gpgsv)
			}
		})
	}
}

var gphdttests = []struct {
	name string
	raw  string
	err  string
	msg  GPHDT
}{
	{
		name: "good sentence",
		raw:  "$GPHDT,123.456,T*32",
		msg: GPHDT{
			Heading: Float64{Valid: true, Value: 123.456},
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

func TestGPHDT(t *testing.T) {
	for _, tt := range gphdttests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gphdt := m.(GPHDT)
				gphdt.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gphdt)
			}
		})
	}
}

var gprmctests = []struct {
	name string
	raw  string
	err  string
	msg  GPRMC
}{
	{
		name: "good sentence A",
		raw:  "$GPRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*70",
		msg: GPRMC{
			Time:      Time{true, 22, 5, 16, 0},
			Validity:  "A",
			Speed:     Float64{Valid: true, Value: 173.8},
			Course:    Float64{Valid: true, Value: 231.8},
			Date:      Date{true, 13, 6, 94},
			Variation: Float64{Valid: true, Value: -4.2},
			Latitude:  MustParseGPS("5133.82 N"),
			Longitude: MustParseGPS("00042.24 W"),
		},
	},
	{
		name: "good sentence B",
		raw:  "$GPRMC,142754.0,A,4302.539570,N,07920.379823,W,0.0,,070617,0.0,E,A*3F",
		msg: GPRMC{
			Time:      Time{true, 14, 27, 54, 0},
			Validity:  "A",
			Speed:     Float64{Valid: true, Value: 0},
			Course:    Float64{Valid: false, Value: 0},
			Date:      Date{true, 7, 6, 17},
			Variation: Float64{Valid: true, Value: 0},
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

func TestGPRMC(t *testing.T) {
	for _, tt := range gprmctests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gprmc := m.(GPRMC)
				gprmc.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gprmc)
			}
		})
	}
}

var gpvtgtests = []struct {
	name string
	raw  string
	err  string
	msg  GPVTG
}{
	{
		name: "good sentence",
		raw:  "$GPVTG,45.5,T,67.5,M,30.45,N,56.40,K*4B",
		msg: GPVTG{
			TrueTrack:        Float64{Valid: true, Value: 45.5},
			MagneticTrack:    Float64{Valid: true, Value: 67.5},
			GroundSpeedKnots: Float64{Valid: true, Value: 30.45},
			GroundSpeedKPH:   Float64{Valid: true, Value: 56.4},
		},
	},
	{
		name: "bad true track",
		raw:  "$GPVTG,T,45.5,67.5,M,30.45,N,56.40,K*4B",
		err:  "nmea: GPVTG invalid true track: T",
	},
}

func TestGPVTG(t *testing.T) {
	for _, tt := range gpvtgtests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gpvtg := m.(GPVTG)
				gpvtg.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gpvtg)
			}
		})
	}
}

var gpzdatests = []struct {
	name string
	raw  string
	err  string
	msg  GPZDA
}{
	{
		name: "good sentence",
		raw:  "$GPZDA,172809.456,12,07,1996,00,00*57",
		msg: GPZDA{
			Time: Time{
				Valid:       true,
				Hour:        17,
				Minute:      28,
				Second:      9,
				Millisecond: 456,
			},
			Day:           Int64{Valid: true, Value: 12},
			Month:         Int64{Valid: true, Value: 7},
			Year:          Int64{Valid: true, Value: 1996},
			OffsetHours:   Int64{Valid: true, Value: 0},
			OffsetMinutes: Int64{Valid: true, Value: 0},
		},
	},
	{
		name: "invalid day",
		raw:  "$GPZDA,220516,D,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*76",
		err:  "nmea: GPZDA invalid day: D",
	},
}

func TestGPZDA(t *testing.T) {
	for _, tt := range gpzdatests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gpzda := m.(GPZDA)
				gpzda.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gpzda)
			}
		})
	}
}
