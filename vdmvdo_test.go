package nmea_test

import (
	"testing"

	"github.com/BertoldVdb/go-ais"
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			NumFragments:   NewInt64(1),
			FragmentNumber: NewInt64(1),
			MessageID:      Int64{},
			Channel:        "A",
			Payload:        []byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0},
			Packet:         ais.PositionReport{Header: ais.Header{MessageID: 0x1, RepeatIndicator: 0x0, UserID: 0xe95fc02}, Valid: true, NavigationalStatus: 0x0, RateOfTurn: -128, Sog: 5, PositionAccuracy: true, Longitude: 6.608731666666667, Latitude: 51.56676166666667, Cog: 113, TrueHeading: 0x1ff, Timestamp: 0x37, SpecialManoeuvreIndicator: 0x1, Spare: 0x0, Raim: true, CommunicationStateNoItdma: ais.CommunicationStateNoItdma{CommunicationState: 0xea0c}},
		},
	},
	{
		name: "Good single fragment message with padding",
		raw:  "!AIVDM,1,1,,A,H77nSfPh4U=<E`H4U8G;:222220,2*1F",
		msg: VDMVDO{
			NumFragments:   NewInt64(1),
			FragmentNumber: NewInt64(1),
			MessageID:      Int64{},
			Channel:        "A",
			Payload:        []byte{0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
			Packet:         ais.StaticDataReport{Header: ais.Header{MessageID: 0x18, RepeatIndicator: 0x0, UserID: 0x1c7da3ba}, Valid: true, Reserved: 0x0, PartNumber: false, ReportA: ais.StaticDataReportA{Valid: true, Name: "LAISSEZFAIRE22"}, ReportB: ais.StaticDataReportB{Valid: false, ShipType: 0x0, VendorIDName: "", VenderIDModel: 0x0, VenderIDSerial: 0x0, CallSign: "", Dimension: ais.FieldDimension{A: 0x0, B: 0x0, C: 0x0, D: 0x0}, FixType: 0x0, Spare: 0x0}},
		},
	},
	{
		name: "Good multipart fragment",
		raw:  "!AIVDM,2,2,4,B,00000000000,2*23",
		msg: VDMVDO{
			NumFragments:   NewInt64(2),
			FragmentNumber: NewInt64(2),
			MessageID:      NewInt64(4),
			Channel:        "B",
			Payload:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	},
	{
		name: "Empty payload",
		raw:  "!AIVDM,1,1,,1,,0*56",
		msg: VDMVDO{
			NumFragments:   NewInt64(1),
			FragmentNumber: NewInt64(1),
			MessageID:      Int64{},
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

var _ = Describe("VDMVDO", func() {
	var (
		parsed VDMVDO
	)
	Describe("Getting data from a $__VDMVDO sentence", func() {
		JustBeforeEach(func() {
			sentence := "!AIVDM,1,1,,B,139`4`0P00PF1l0MUSjN4?vJ2L5H,0*6A"
			parseResult, err := Parse(sentence)
			if err != nil {
				Fail("Could not parse sentence")
			}
			var ok bool
			if parsed, ok = parseResult.(VDMVDO); !ok {
				Fail("Could not cast to VDMVDO")
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid MMSI", func() {
				Expect(parsed.GetMMSI()).To(Equal(uint32(211420320)))
			})
			It("should give a valid Position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(Equal(51.70678833333333))
				Expect(lon).To(Equal(4.81216))
			})
		})
	})
})
