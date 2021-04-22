package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var dpttests = []struct {
	name string
	raw  string
	err  string
	msg  DPT
}{
	{
		name: "good sentence",
		raw:  "$SDDPT,0.5,0.5,*7B",
		msg: DPT{
			Depth:      NewFloat64(0.5),
			Offset:     NewFloat64(0.5),
			RangeScale: Float64{},
		},
	},
	{
		name: "good sentence with scale",
		raw:  "$SDDPT,0.5,0.5,0.1*54",
		msg: DPT{
			Depth:      NewFloat64(0.5),
			Offset:     NewFloat64(0.5),
			RangeScale: NewFloat64(0.1),
		},
	},
	{
		name: "bad validity",
		raw:  "$SDDPT,0.5,0.5,*AA",
		err:  "nmea: sentence checksum mismatch [7B != AA]",
	},
}

func TestDPT(t *testing.T) {
	for _, tt := range dpttests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				dpt := m.(DPT)
				dpt.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, dpt)
			}
		})
	}
}

var _ = Describe("DPT", func() {
	var (
		parsed DPT
	)
	Describe("Getting data from a $__DPT sentence", func() {
		BeforeEach(func() {
			parsed = DPT{
				Depth: NewFloat64(DepthBelowSurfaceMeters - DepthTransducerMeters),
			}
		})
		Context("When having a parsed sentence and a positive offset", func() {
			JustBeforeEach(func() {
				parsed.Offset = NewFloat64(DepthTransducerMeters)
			})
			It("should give a valid depth below transducer", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(Float64Equal(DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
			It("should give a valid depth below surface", func() {
				Expect(parsed.GetDepthBelowSurface()).To(Float64Equal(DepthBelowSurfaceMeters, 0.00001))
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowKeel()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When having a parsed sentence and a negative offset", func() {
			JustBeforeEach(func() {
				parsed.Offset = NewFloat64(DepthTransducerMeters - DepthKeelMeters)
			})
			It("should give a valid depth below transducer", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(Float64Equal(DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid depth below keel", func() {
				Expect(parsed.GetDepthBelowKeel()).To(Float64Equal(DepthBelowSurfaceMeters-DepthKeelMeters, 0.00001))
			})
		})
		Context("When having a parsed sentence and no offset", func() {
			JustBeforeEach(func() {
				parsed.Offset = Float64{}
			})
			It("should give a valid depth below transducer", func() {
				Expect(parsed.GetDepthBelowTransducer()).To(Float64Equal(DepthBelowSurfaceMeters-DepthTransducerMeters, 0.00001))
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowKeel()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When having a parsed sentence and no depth", func() {
			JustBeforeEach(func() {
				parsed.Depth = Float64{}
				parsed.Offset = NewFloat64(DepthTransducerMeters)
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowTransducer()
				Expect(err).To(HaveOccurred())
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowSurface()
				Expect(err).To(HaveOccurred())
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetDepthBelowKeel()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
