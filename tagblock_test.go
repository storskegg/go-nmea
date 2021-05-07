package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tagblock", func() {
	Describe("Testing Parse method", func() {
		Context("with a basic tag block", func() {
			It("returns a valid value", func() {
				result, _ := Parse("\\s:Satelite_1,c:1553390539*62\\$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				rmc := result.(RMC)
				Expect(rmc.TagBlock.Valid).To(BeTrue())
				Expect(rmc.TagBlock.Time).To(Equal(NewInt64(1553390539)))
				Expect(rmc.TagBlock.Source).To(Equal(NewString("Satelite_1")))
			})
		})
		Context("with a tag block with an unknown tag", func() {
			It("returns a valid value", func() {
				result, _ := Parse("\\x:NorSat_1,c:1564827317*42\\$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				rmc := result.(RMC)
				Expect(rmc.TagBlock.Valid).To(BeTrue())
				Expect(rmc.TagBlock.Time).To(Equal(NewInt64(1564827317)))
				Expect(rmc.TagBlock.Source).To(Equal(NewInvalidString("not specified")))
			})
		})
		Context("with a tag block with an unknown tag and ten millisecond timestamp", func() {
			It("returns a valid value", func() {
				result, _ := Parse("\\x:NorSat_1,c:1564827317000*72\\$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				rmc := result.(RMC)
				Expect(rmc.TagBlock.Valid).To(BeTrue())
				Expect(rmc.TagBlock.Time).To(Equal(NewInt64(1564827317000)))
				Expect(rmc.TagBlock.Source).To(Equal(NewInvalidString("not specified")))
			})
		})
		Context("with a tag block with all tags", func() {
			It("returns a valid value", func() {
				result, _ := Parse("\\s:satelite,c:1564827317,r:1553390539,d:ara,g:bulk,n:13,t:helloworld*3F\\$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				rmc := result.(RMC)
				Expect(rmc.TagBlock.Valid).To(BeTrue())
				Expect(rmc.TagBlock.Time).To(Equal(NewInt64(1564827317)))
				Expect(rmc.TagBlock.RelativeTime).To(Equal(NewInt64(1553390539)))
				Expect(rmc.TagBlock.Destination).To(Equal(NewString("ara")))
				Expect(rmc.TagBlock.Grouping).To(Equal(NewString("bulk")))
				Expect(rmc.TagBlock.Source).To(Equal(NewString("satelite")))
				Expect(rmc.TagBlock.Text).To(Equal(NewString("helloworld")))
				Expect(rmc.TagBlock.LineCount).To(Equal(NewInt64(13)))
			})
		})
		Context("with a tag block with an empty tag", func() {
			It("returns a valid value", func() {
				result, _ := Parse("\\s:satelite,,c:1564827317,r:1553390539,d:ara,g:bulk,n:13,t:helloworld*13\\$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				rmc := result.(RMC)
				Expect(rmc.TagBlock.Valid).To(BeFalse())
				Expect(rmc.TagBlock.InvalidReason).To(Equal("nmea: tagblock field is malformed (should be <key>:<value>) []"))
			})
		})
		Context("with a tag block with an invalid checksum", func() {
			It("returns a valid value", func() {
				result, _ := Parse("\\s:satelite,c:1564827317,r:1553390539,d:ara,g:bulk,n:13,t:helloworld*3A\\$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				rmc := result.(RMC)
				Expect(rmc.TagBlock.Valid).To(BeFalse())
				Expect(rmc.TagBlock.InvalidReason).To(Equal("nmea: tagblock checksum mismatch [3F != 3A]"))
			})
		})
		Context("with a tag block without a checksum", func() {
			It("returns a valid value", func() {
				result, _ := Parse("\\s:satelite,c:156482731749\\$GPRMC,001225,A,2832.1834,N,08101.0536,W,12,25,251211,1.2,E,A*03")
				rmc := result.(RMC)
				Expect(rmc.TagBlock.Valid).To(BeFalse())
				Expect(rmc.TagBlock.InvalidReason).To(Equal("nmea: tagblock does not contain checksum separator"))
			})
		})
	})
})
