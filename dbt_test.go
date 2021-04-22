package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var dbttests = []struct {
	name string
	raw  string
	err  string
	msg  DBT
}{
	{
		name: "good sentence",
		raw:  "$IIDBT,032.93,f,010.04,M,005.42,F*2C",
		msg: DBT{
			DepthFeet:    MustParseDecimal("32.93"),
			DepthMeters:  MustParseDecimal("10.04"),
			DepthFathoms: MustParseDecimal("5.42"),
		},
	},
	{
		name: "bad validity",
		raw:  "$IIDBT,032.93,f,010.04,M,005.42,F*22",
		err:  "nmea: sentence checksum mismatch [2C != 22]",
	},
}

func TestDBT(t *testing.T) {
	for _, tt := range dbttests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				dbt := m.(DBT)
				dbt.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, dbt)
			}
		})
	}
}

var _ = Describe("DBT", func() {
	var (
		parsed DBT
	)
	Describe("Getting data from a $__DBT sentence", func() {
		BeforeEach(func() {
			parsed = DBT{
				DepthFeet:    NewFloat64(DepthBelowSurfaceFeet - DepthTransducerFeet),
				DepthMeters:  NewFloat64(DepthBelowSurfaceMeters - DepthTransducerMeters),
				DepthFathoms: NewFloat64(DepthBelowSurfaceFathoms - DepthTransducerFanthoms),
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(Float64Equal(DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with only depth in feet set", func() {
			JustBeforeEach(func() {
				parsed.DepthMeters = Float64{}
				parsed.DepthFathoms = Float64{}
			})
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(Float64Equal(DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with only depth in fathoms set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = Float64{}
				parsed.DepthMeters = Float64{}
			})
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(Float64Equal(DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with only depth in meters set", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = Float64{}
				parsed.DepthFathoms = Float64{}
			})
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(Float64Equal(DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing depth values", func() {
			JustBeforeEach(func() {
				parsed.DepthFeet = Float64{}
				parsed.DepthMeters = Float64{}
				parsed.DepthFathoms = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowTransducer()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
