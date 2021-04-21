package nmea

import (
	"testing"

	"github.com/BertoldVdb/go-ais"
	"github.com/stretchr/testify/assert"
)

var vdmtests = []struct {
	name string
	raw  string
	err  string
	msg  VDMVDO
}{
	{
		name: "Good single fragment message",
		raw:  "!AIVDM,1,1,,A,13aGt0PP0jPN@9fMPKVDJgwfR>`<,0*55",
		msg: VDMVDO{
			NumFragments:   Int64{Valid: true, Value: 1},
			FragmentNumber: Int64{Valid: true, Value: 1},
			MessageID:      Int64{Valid: false, Value: 0},
			Channel:        "A",
			Payload:        []byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0},
			Packet:         ais.PositionReport{Header: ais.Header{MessageID: 0x1, RepeatIndicator: 0x0, UserID: 0xe95fc02}, Valid: true, NavigationalStatus: 0x0, RateOfTurn: -128, Sog: 5, PositionAccuracy: true, Longitude: 6.608731666666667, Latitude: 51.56676166666667, Cog: 113, TrueHeading: 0x1ff, Timestamp: 0x37, SpecialManoeuvreIndicator: 0x1, Spare: 0x0, Raim: true, CommunicationStateNoItdma: ais.CommunicationStateNoItdma{CommunicationState: 0xea0c}},
		},
	},
	{
		name: "Good single fragment message with padding",
		raw:  "!AIVDM,1,1,,A,H77nSfPh4U=<E`H4U8G;:222220,2*1F",
		msg: VDMVDO{
			NumFragments:   Int64{Valid: true, Value: 1},
			FragmentNumber: Int64{Valid: true, Value: 1},
			MessageID:      Int64{Valid: false, Value: 0},
			Channel:        "A",
			Payload:        []byte{0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
			Packet:         ais.StaticDataReport{Header: ais.Header{MessageID: 0x18, RepeatIndicator: 0x0, UserID: 0x1c7da3ba}, Valid: true, Reserved: 0x0, PartNumber: false, ReportA: ais.StaticDataReportA{Valid: true, Name: "LAISSEZFAIRE22"}, ReportB: ais.StaticDataReportB{Valid: false, ShipType: 0x0, VendorIDName: "", VenderIDModel: 0x0, VenderIDSerial: 0x0, CallSign: "", Dimension: ais.FieldDimension{A: 0x0, B: 0x0, C: 0x0, D: 0x0}, FixType: 0x0, Spare: 0x0}},
		},
	},
	{
		name: "Good multipart fragment",
		raw:  "!AIVDM,2,2,4,B,00000000000,2*23",
		msg: VDMVDO{
			NumFragments:   Int64{Valid: true, Value: 2},
			FragmentNumber: Int64{Valid: true, Value: 2},
			MessageID:      Int64{Valid: true, Value: 4},
			Channel:        "B",
			Payload:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	},
	{
		name: "Empty payload",
		raw:  "!AIVDM,1,1,,1,,0*56",
		msg: VDMVDO{
			NumFragments:   Int64{Valid: true, Value: 1},
			FragmentNumber: Int64{Valid: true, Value: 1},
			MessageID:      Int64{Valid: false, Value: 0},
			Channel:        "1",
			Payload:        []byte{},
		},
	},
	{
		name: "Invalid number of fragments",
		raw:  "!AIVDM,x,1,,1,000 00,0*0F",
		err:  "nmea: AIVDM invalid number of fragments: x",
	},
	{
		name: "Invalid symbol in payload",
		raw:  "!AIVDM,1,1,,1,000 00,0*46",
		err:  "nmea: AIVDM invalid payload: data byte",
	},
	{
		name: "Negative number of fill bits",
		raw:  "!AIVDM,1,1,,1,000,-3*48",
		err:  "nmea: AIVDM invalid payload: fill bits",
	},
	{
		name: "Too high number of fill bits",
		raw:  "!AIVDO,1,1,,1,000,20*56",
		err:  "nmea: AIVDO invalid payload: fill bits",
	},
	{
		name: "Negative number of bits",
		raw:  "!AIVDM,1,1,,1,,2*54",
		err:  "nmea: AIVDM invalid payload: num bits",
	},
}

func TestVDM(t *testing.T) {
	for _, tt := range vdmtests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)

			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				vdm := m.(VDMVDO)
				vdm.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, vdm)
			}
		})
	}
}
