package nmea_test

import (
	"time"

	"github.com/BertoldVdb/go-ais"
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("VDMVDO", func() {
	var (
		sentence Sentence
		parsed   VDMVDO
		err      error
		raws     []string
	)
	JustBeforeEach(func() {
		for _, raw := range raws {
			sentence, err = Parse(raw)
			if err != nil {
				break
			}
		}
		if sentence != nil {
			parsed = sentence.(VDMVDO)
		} else {
			parsed = VDMVDO{}
		}
	})
	Describe("Parsing", func() {
		Context("a valid single fragment sentence", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,A,13aGt0PP0jPN@9fMPKVDJgwfR>`<,0*55",
				}
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
				raws = []string{
					"!AIVDM,1,1,,A,H77nSfPh4U=<E`H4U8G;:222220,2*1F",
				}
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
				raws = []string{
					"!AIVDM,2,2,4,B,00000000000,2*23",
				}
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
				raws = []string{
					"!AIVDM,1,1,,1,,0*56",
				}
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
				raws = []string{
					"!AIVDM,x,1,,1,000 00,0*0F",
				}
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid number of fragments: x"))
			})
		})
		Context("a sentence with an invalid symbol in the payload", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,1,000 00,0*46",
				}
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid payload: data byte"))
			})
		})
		Context("a sentence with a negative number of fill bits", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,1,000,-3*48",
				}
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid payload: fill bits"))
			})
		})
		Context("a sentence with too much fill bits", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDO,1,1,,1,000,20*56",
				}
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDO invalid payload: fill bits"))
			})
		})
		Context("a sentence with a negative number of bits", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,1,,2*54",
				}
			})
			It("returns an errors", func() {
				Expect(err).To(MatchError("nmea: AIVDM invalid payload: num bits"))
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,A,13aGt0PP0jPN@9fMPKVDJgwfR>`<,0*54",
				}
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [55 != 54]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
		Context("when having a scheduled position report", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,B,13aL>lwP0rPF<=8MSWjWSwwH2<3d,0*24",
				}
			})
			It("returns a valid MMSI", func() {
				Expect(parsed.GetMMSI()).To(Equal("244780755"))
			})
			It("returns an error", func() {
				_, err := parsed.GetCallSign()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetENINumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetIMONumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				Expect(parsed.GetNavigationStatus()).To(Equal("default"))
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselBeam()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselLength()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselName()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselType()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetRateOfTurn()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", 3.3772121026090276, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(Equal(51.65388333333333))
				Expect(lon).To(Equal(4.8476333333333335))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", 2.9837752, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetDestination()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a scheduled position report", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,A,13u?etPv2;0n:dDPwUM1U1Cb069D,0*24",
				}
			})
			It("returns a valid MMSI", func() {
				Expect(parsed.GetMMSI()).To(Equal("265547250"))
			})
			It("returns an error", func() {
				_, err := parsed.GetCallSign()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetENINumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetIMONumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid navigation status", func() {
				Expect(parsed.GetNavigationStatus()).To(Equal("motoring"))
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselBeam()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselLength()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselName()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselType()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid rate of turn", func() {
				Expect(parsed.GetRateOfTurn()).To(BeNumerically("~", -0.050614548308/60.0, 0.0001))
			})
			It("returns a valid course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", 0.70511301781, 0.00001))
			})
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(BeNumerically("~", 0.71558499332, 0.00001))
			})
			It("returns a valid position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(BeNumerically("~", 57.6603533, 0.00001))
				Expect(lon).To(BeNumerically("~", 11.8329767, 0.00001))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", 7.1507777778, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetDestination()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having another scheduled position report", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,1,1,,A,13aHE1PPOwPGJSJMiB8N4;=4P<23,0*7E",
				}
			})
			It("returns a valid MMSI", func() {
				Expect(parsed.GetMMSI()).To(Equal("244716806"))
			})
			It("returns an error", func() {
				_, err := parsed.GetCallSign()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetENINumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetIMONumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid navigation status", func() {
				Expect(parsed.GetNavigationStatus()).To(Equal("motoring"))
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselBeam()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselLength()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselName()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselType()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid rate of turn", func() {
				Expect(parsed.GetRateOfTurn()).To(BeNumerically("~", -0.2120575, 0.0001))
			})
			It("returns a valid course over ground", func() {
				_, err := parsed.GetTrueCourseOverGround()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(BeNumerically("~", 6.248279, 0.00001))
			})
			It("returns a valid position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(BeNumerically("~", 52.026935, 0.00001))
				Expect(lon).To(BeNumerically("~", 5.11506166666667, 0.00001))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", 52.6276212, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetDestination()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when having a ship static and voyage related data report", func() {
			BeforeEach(func() {
				raws = []string{
					"!AIVDM,2,1,1,B,53aDr?H000010CS7OH04@Dh4q@D000000000001?1QR75u8kP05iDRiC,0*11",
					"!AIVDM,2,2,1,B,Q0C@00000000000,2*44",
				}
			})
			It("returns a valid MMSI", func() {
				Expect(parsed.GetMMSI()).To(Equal("244660797"))
			})
			It("returns a valid call sign", func() {
				Expect(parsed.GetCallSign()).To(Equal("PD8176"))
			})
			It("returns an error", func() {
				_, err := parsed.GetENINumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid IMO number", func() {
				Expect(parsed.GetIMONumber()).To(Equal("0"))
			})
			It("returns an error", func() {
				_, err := parsed.GetNavigationStatus()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid beam", func() {
				Expect(parsed.GetVesselBeam()).To(BeNumerically("~", 12, 0.00001))
			})
			It("returns a valid length", func() {
				Expect(parsed.GetVesselLength()).To(BeNumerically("~", 110, 0.00001))
			})
			It("returns a valid vessel name", func() {
				Expect(parsed.GetVesselName()).To(Equal("ADELANTE"))
			})
			It("returns a valid vessel type", func() {
				Expect(parsed.GetVesselType()).To(Equal("Cargo, No additional information"))
			})
			It("returns an error", func() {
				_, err := parsed.GetRateOfTurn()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueCourseOverGround()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, _, err := parsed.GetPosition2D()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetSpeedOverGround()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid destination", func() {
				Expect(parsed.GetDestination()).To(Equal("WERKENDAM"))
			})
			It("returns a valid eta", func() {
				eta, _ := parsed.GetETA()
				// Can't validate year because it depends on the actual date when the test is run
				Expect(eta.Month()).To(Equal(time.Month(4)))
				Expect(eta.Day()).To(Equal(17))
				Expect(eta.Hour()).To(Equal(19))
				Expect(eta.Minute()).To(Equal(32))
				Expect(eta.Second()).To(Equal(0))
				Expect(eta.Nanosecond()).To(Equal(0))
				Expect(eta.Zone()).To(Equal("UTC"))
			})
		})
	})
	Describe("Getting data from a VDMVDO struct", func() {
		Context("with a position report", func() {
			var (
				parsed VDMVDO
			)
			BeforeEach(func() {
				parsed = VDMVDO{
					Packet: ais.PositionReport{
						Header: ais.Header{
							UserID: 244660797,
						},
						Valid:                     true,
						NavigationalStatus:        0,
						RateOfTurn:                int16(-40),
						Sog:                       ais.Field10(SpeedOverGroundKnots),
						PositionAccuracy:          false,
						Longitude:                 ais.FieldLatLonFine(Longitude),
						Latitude:                  ais.FieldLatLonFine(Latitude),
						Cog:                       ais.Field10(TrueDirectionDegrees),
						TrueHeading:               uint16(212),
						Timestamp:                 0,
						SpecialManoeuvreIndicator: 0,
						Spare:                     0,
						Raim:                      false,
					},
				}
			})
			It("returns a valid MMSI", func() {
				Expect(parsed.GetMMSI()).To(Equal("244660797"))
			})
			It("returns an error", func() {
				_, err := parsed.GetCallSign()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetENINumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetIMONumber()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				Expect(parsed.GetNavigationStatus()).To(Equal("motoring"))
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselBeam()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselLength()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselName()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetVesselType()
				Expect(err).To(HaveOccurred())
			})
			It("returns a valid rate of turn", func() {
				Expect(parsed.GetRateOfTurn()).To(BeNumerically("~", -0.020776533612, 0.00001))
			})
			It("returns an error", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(BeNumerically("~", TrueDirectionRadians, 0.00001))
			})
			It("returns a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(BeNumerically("~", 3.7000980142, 0.00001))
			})
			It("returns a valid position", func() {
				lat, lon, _ := parsed.GetPosition2D()
				Expect(lat).To(Equal(Latitude))
				Expect(lon).To(Equal(Longitude))
			})
			It("returns a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(BeNumerically("~", SpeedOverGroundMPS, 0.00001))
			})
			It("returns an error", func() {
				_, err := parsed.GetDestination()
				Expect(err).To(HaveOccurred())
			})
			It("returns an error", func() {
				_, err := parsed.GetETA()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
