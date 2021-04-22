package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var vhw = []struct {
	name string
	raw  string
	err  string
	msg  VHW
}{
	{
		name: "good sentence",
		raw:  "$VWVHW,45.0,T,43.0,M,3.5,N,6.4,K*56",
		msg: VHW{
			TrueHeading:            NewFloat64(45.0),
			MagneticHeading:        NewFloat64(43.0),
			SpeedThroughWaterKnots: NewFloat64(3.5),
			SpeedThroughWaterKPH:   NewFloat64(6.4),
		},
	},
	{
		name: "bad sentence",
		raw:  "$VWVHW,T,45.0,43.0,M,3.5,N,6.4,K*56",
		err:  "nmea: VWVHW invalid true heading: T",
	},
}

func TestVHW(t *testing.T) {
	for _, tt := range vhw {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				vhw := m.(VHW)
				vhw.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, vhw)
			}
		})
	}
}

var _ = Describe("VHW", func() {
	var (
		parsed VHW
	)
	Describe("Getting data from a $__VHW sentence", func() {
		BeforeEach(func() {
			parsed = VHW{
				TrueHeading:            NewFloat64(TrueDirectionDegrees),
				MagneticHeading:        NewFloat64(MagneticDirectionDegrees),
				SpeedThroughWaterKPH:   NewFloat64(SpeedThroughWaterKPH),
				SpeedThroughWaterKnots: NewFloat64(SpeedThroughWaterKnots),
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing true heading", func() {
			JustBeforeEach(func() {
				parsed.TrueHeading = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing magnetic track", func() {
			JustBeforeEach(func() {
				parsed.MagneticHeading = Float64{}
			})
			It("should give a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetMagneticHeading()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing speed over ground kph", func() {
			JustBeforeEach(func() {
				parsed.SpeedThroughWaterKPH = Float64{}
			})
			It("should give a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing speed over ground knots", func() {
			JustBeforeEach(func() {
				parsed.SpeedThroughWaterKnots = Float64{}
			})
			It("should give a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed through water", func() {
				Expect(parsed.GetSpeedThroughWater()).To(Float64Equal(SpeedThroughWaterMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing speed over ground kph and knots", func() {
			JustBeforeEach(func() {
				parsed.SpeedThroughWaterKPH = Float64{}
				parsed.SpeedThroughWaterKnots = Float64{}
			})
			It("should give a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic heading", func() {
				Expect(parsed.GetMagneticHeading()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetSpeedThroughWater()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
