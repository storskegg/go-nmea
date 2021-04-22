package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var dbstests = []struct {
	name string
	raw  string
	err  string
	msg  DBS
}{
	{
		name: "good sentence",
		raw:  "$23DBS,01.9,f,0.58,M,00.3,F*21",
		msg: DBS{
			DepthFeet:    MustParseDecimal("1.9"),
			DepthMeters:  MustParseDecimal("0.58"),
			DepthFathoms: MustParseDecimal("0.3"),
		},
	},
	{
		name: "bad validity",
		raw:  "$23DBS,01.9,f,0.58,M,00.3,F*25",
		err:  "nmea: sentence checksum mismatch [21 != 25]",
	},
}

func TestDBS(t *testing.T) {
	for _, tt := range dbstests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				dbs := m.(DBS)
				dbs.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, dbs)
			}
		})
	}
}

var _ = Describe("DBS", func() {
	var (
		parsed DBS
	)
	Describe("Getting data from a $__DBS sentence", func() {
		BeforeEach(func() {
			parsed = DBS{
				DepthFeet:    NewFloat64(DepthBelowSurfaceFeet),
				DepthMeters:  NewFloat64(DepthBelowSurfaceMeters),
				DepthFathoms: NewFloat64(DepthBelowSurfaceFathoms),
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(Float64Equal(DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with only depth in feet set", func() {
			JustBeforeEach(func() {
				parsed.DepthMeters = Float64{}
				parsed.DepthFathoms = Float64{}
			})
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(Float64Equal(DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with only depth in fathoms set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = Float64{}
				parsed.DepthMeters = Float64{}
			})
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(Float64Equal(DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with only depth in meters set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = Float64{}
				parsed.DepthFathoms = Float64{}
			})
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(Float64Equal(DepthBelowSurfaceMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing depth values", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = Float64{}
				parsed.DepthMeters = Float64{}
				parsed.DepthFathoms = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
