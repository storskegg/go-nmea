package nmea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gnstests = []struct {
	name string
	raw  string
	err  string
	msg  GNS
}{
	{
		name: "good sentence A",
		raw:  "$GNGNS,014035.00,4332.69262,S,17235.48549,E,RR,13,0.9,25.63,11.24,,*70",
		msg: GNS{
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
		msg: GNS{
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
		msg: GNS{
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

func TestGNS(t *testing.T) {
	for _, tt := range gnstests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gns := m.(GNS)
				gns.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, gns)
			}
		})
	}
}
