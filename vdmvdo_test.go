package nmea_test

import (
	"github.com/BertoldVdb/go-ais"
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("VDMVDO", func() {
	var (
		sentence Sentence
		parsed   VDMVDO
		err      error
		raw      string
	)
	Describe("Parsing", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
			if sentence != nil {
				parsed = sentence.(VDMVDO)
			} else {
				parsed = VDMVDO{}
			}
		})
		Context("a valid single fragment sentence", func() {
			BeforeEach(func() {
				raw = "!AIVDM,1,1,,A,13aGt0PP0jPN@9fMPKVDJgwfR>`<,0*55"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid VDMVDO struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumFragments":   Equal(NewInt64(1)),
					"FragmentNumber": Equal(NewInt64(1)),
					"MessageID":      Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
					"Channel":        Equal(NewString("A")),
					"Payload":        Equal([]byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0}),
					"Packet":         Equal(ais.PositionReport{Header: ais.Header{MessageID: 0x1, RepeatIndicator: 0x0, UserID: 0xe95fc02}, Valid: true, NavigationalStatus: 0x0, RateOfTurn: -128, Sog: 5, PositionAccuracy: true, Longitude: 6.608731666666667, Latitude: 51.56676166666667, Cog: 113, TrueHeading: 0x1ff, Timestamp: 0x37, SpecialManoeuvreIndicator: 0x1, Spare: 0x0, Raim: true, CommunicationStateNoItdma: ais.CommunicationStateNoItdma{CommunicationState: 0xea0c}}),
				}))
			})
		})
		Context("a valid single fragment sentence with padding", func() {
			BeforeEach(func() {
				raw = "!AIVDM,1,1,,A,H77nSfPh4U=<E`H4U8G;:222220,2*1F"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid VDMVDO struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumFragments":   Equal(NewInt64(1)),
					"FragmentNumber": Equal(NewInt64(1)),
					"MessageID":      Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
					"Channel":        Equal(NewString("A")),
					"Payload":        Equal([]byte{0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}),
					"Packet":         Equal(ais.StaticDataReport{Header: ais.Header{MessageID: 0x18, RepeatIndicator: 0x0, UserID: 0x1c7da3ba}, Valid: true, Reserved: 0x0, PartNumber: false, ReportA: ais.StaticDataReportA{Valid: true, Name: "LAISSEZFAIRE22"}, ReportB: ais.StaticDataReportB{Valid: false, ShipType: 0x0, VendorIDName: "", VenderIDModel: 0x0, VenderIDSerial: 0x0, CallSign: "", Dimension: ais.FieldDimension{A: 0x0, B: 0x0, C: 0x0, D: 0x0}, FixType: 0x0, Spare: 0x0}}),
				}))
			})
		})
		Context("a valid multipart sentence", func() {
			BeforeEach(func() {
				raw = "!AIVDM,2,2,4,B,00000000000,2*23"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid VDMVDO struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumFragments":   Equal(NewInt64(2)),
					"FragmentNumber": Equal(NewInt64(2)),
					"MessageID":      Equal(NewInt64(4)),
					"Channel":        Equal(NewString("B")),
					"Payload":        Equal([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
					"Packet":         BeNil(),
				}))
			})
		})
		Context("a valid  sentence with an empty payload", func() {
			BeforeEach(func() {
				raw = "!AIVDM,1,1,,1,,0*56"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("equals a valid VDMVDO struct", func() {
				Expect(parsed).To(MatchFields(IgnoreExtras, Fields{
					"NumFragments":   Equal(NewInt64(1)),
					"FragmentNumber": Equal(NewInt64(1)),
					"MessageID":      Equal(NewInvalidInt64("strconv.ParseInt: parsing \"\": invalid syntax")),
					"Channel":        Equal(NewString("1")),
					"Payload":        Equal([]byte{}),
					"Packet":         BeNil(),
				}))
			})
		})
		Context("a sentence with an invalid number of fragments", func() {
			BeforeEach(func() {
				raw = "!AIVDM,x,1,,1,000 00,0*0F"
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid number of fragments: x"))
			})
		})
		Context("a sentence with an invalid symbol in the payload", func() {
			BeforeEach(func() {
				raw = "!AIVDM,1,1,,1,000 00,0*46"
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid payload: data byte"))
			})
		})
		Context("a sentence with a negative number of fill bits", func() {
			BeforeEach(func() {
				raw = "!AIVDM,1,1,,1,000,-3*48"
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid payload: fill bits"))
			})
		})
		Context("a sentence with too much fill bits", func() {
			BeforeEach(func() {
				raw = "!AIVDO,1,1,,1,000,20*56"
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDO invalid payload: fill bits"))
			})
		})
		Context("a sentence with a negative number of bits", func() {
			BeforeEach(func() {
				raw = "!AIVDM,1,1,,1,,2*54"
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid payload: num bits"))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "!AIVDM,1,1,,A,13aGt0PP0jPN@9fMPKVDJgwfR>`<,0*54"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [55 != 54]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Getting data from a VDMVDO struct", func() {
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
		Context("when having a complete struct", func() {
			It("returns a valid MMSI", func() {
				Expect(parsed.GetMMSI()).To(Equal(uint32(211420320)))
			})
			It("returns a valid Position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(Equal(51.70678833333333))
				Expect(lon).To(Equal(4.81216))
			})
		})
	})
})
