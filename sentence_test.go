package nmea_test

import (
	"errors"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sentence", func() {
	Describe("Testing methods of BaseSentence", func() {
		var (
			bs BaseSentence
		)
		BeforeEach(func() {
			bs = BaseSentence{
				Talker:   "GN",
				Type:     "RMC",
				Fields:   []string{"001225", "A", "2832.1834", "N", "08101.0536", "W", "12", "25", "251211", "1.2", "E", "A"},
				Checksum: "03",
				Raw:      "$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03",
				TagBlock: TagBlock{},
			}
		})
		Context("when getting the prefix", func() {
			It("returns a valid value", func() {
				Expect(bs.Prefix()).To(Equal("GNRMC"))
			})
		})
		Context("when getting the datatype", func() {
			It("returns a valid value", func() {
				Expect(bs.DataType()).To(Equal("RMC"))
			})
		})
		Context("when getting the talker id", func() {
			It("returns a valid value", func() {
				Expect(bs.TalkerID()).To(Equal("GN"))
			})
		})
		Context("when getting the string representation", func() {
			It("returns a valid value", func() {
				Expect(bs.String()).To(Equal("$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03"))
			})
		})
	})
	Describe("Testing the (Must)RegisterParser function", func() {
		cp := func(BaseSentence) (Sentence, error) {
			return nil, nil
		}
		Context("when a custom parser has not been registered yet", func() {
			It("returns no error", func() {
				Expect(RegisterParser("CPA", cp)).ToNot(HaveOccurred())
			})
		})
		Context("when a custom parser is registered twice", func() {
			It("returns an error", func() {
				Expect(RegisterParser("CPB", cp)).ToNot(HaveOccurred())
				Expect(RegisterParser("CPB", cp)).To(MatchError("nmea: parser for sentence type '\"CPB\"' already exists"))
			})
		})
		Context("when a custom parser has not been registered yet", func() {
			It("returns no error", func() {
				Expect(func() { MustRegisterParser("CPC", cp) }).ToNot(Panic())
			})
		})
		Context("when a custom parser is registered twice", func() {
			It("returns an error", func() {
				Expect(func() { MustRegisterParser("CPD", cp) }).ToNot(Panic())
				Expect(func() { MustRegisterParser("CPD", cp) }).To(PanicWith(errors.New("nmea: parser for sentence type '\"CPD\"' already exists")))
			})
		})
		Context("when registering and using a custom parser", func() {
			It("parses the custom sentence", func() {
				type XYZType struct {
					BaseSentence
					Time    Time
					Counter Int64
					Label   String
					Value   Float64
				}
				err := RegisterParser("XYZ", func(s BaseSentence) (Sentence, error) {
					p := NewParser(s)
					return XYZType{
						BaseSentence: s,
						Time:         p.Time(0, "time"),
						Label:        p.String(1, "label"),
						Counter:      p.Int64(2, "counter"),
						Value:        p.Float64(3, "value"),
					}, p.Err()
				})
				Expect(err).ToNot(HaveOccurred())
				sentence := "$00XYZ,220516,A,23,5133.82,W*42"
				s, err := Parse(sentence)
				Expect(err).ToNot(HaveOccurred())

				_, ok := s.(XYZType)
				Expect(ok).To(BeTrue())
			})
		})
	})
	Describe("Testing the Parse function", func() {
		Context("when a standard sentence is given", func() {
			It("returns a valid value", func() {
				result, err := Parse("$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				Expect(result).ToNot(BeNil())
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("when a unsupported sentence is given", func() {
			It("returns an error", func() {
				result, err := Parse("$PSTIS,*61")
				Expect(result).To(BeNil())
				Expect(err).To(MatchError("nmea: sentence prefix 'PSTIS' not supported"))
			})
		})
		Context("when a standard sentence is given with leading an trailing spaces", func() {
			It("returns a valid value", func() {
				result, err := Parse("     $GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03        				")
				Expect(result).ToNot(BeNil())
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("when a standard sentence is given with a valid TAG block", func() {
			It("returns a valid value", func() {
				result, err := Parse("\\s:Satellite_1,c:1553390539*62\\!AIVDM,1,1,,A,13M@ah0025QdPDTCOl`K6`nV00Sv,0*52")
				Expect(result).ToNot(BeNil())
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("when a standard sentence is given with a bad checksum", func() {
			It("returns an error", func() {
				result, err := Parse("$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*04")
				Expect(result).To(BeNil())
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [03 != 04]"))
			})
		})
		Context("when a standard sentence is given with a bad start character", func() {
			It("returns an error", func() {
				result, err := Parse("%%GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				Expect(result).To(BeNil())
				Expect(err).To(MatchError("nmea: sentence does not start with a '$' or '!'"))
			})
		})
		Context("when a standard sentence is given without a checksum separator", func() {
			It("returns an error", func() {
				result, err := Parse("$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A")
				Expect(result).To(BeNil())
				Expect(err).To(MatchError("nmea: sentence does not contain checksum separator"))
			})
		})
		Context("when a standard sentence is given without a start delimiter", func() {
			It("returns an error", func() {
				result, err := Parse("abc$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				Expect(result).To(BeNil())
				Expect(err).To(MatchError("nmea: sentence does not start with a '$' or '!'"))
			})
		})
		Context("when a standard sentence is given without a TAG Block start delimiter", func() {
			It("returns an error", func() {
				result, err := Parse("s:Satellite_1,c:1553390539*62\\!AIVDM,1,1,,A,13M@ah0025QdPDTCOl`K6`nV00Sv,0*52")
				Expect(result).To(BeNil())
				Expect(err).To(MatchError("nmea: sentence does not start with a '$' or '!'"))
			})
		})
		Context("when a standard sentence is given without a TAG Block end delimiter", func() {
			It("returns an error", func() {
				result, err := Parse("\\s:Satellite_1,c:1553390539*62!AIVDM,1,1,,A,13M@ah0025QdPDTCOl`K6`nV00Sv,0*52")
				Expect(result).To(BeNil())
				Expect(err).To(MatchError("nmea: sentence does not start with a '$' or '!'"))
			})
		})
	})
})
