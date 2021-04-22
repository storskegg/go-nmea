package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var vtgtests = []struct {
	name string
	raw  string
	err  string
	msg  VTG
}{
	{
		name: "good sentence",
		raw:  "$GPVTG,45.5,T,67.5,M,30.45,N,56.40,K*4B",
		msg: VTG{
			TrueTrack:        NewFloat64(45.5),
			MagneticTrack:    NewFloat64(67.5),
			GroundSpeedKnots: NewFloat64(30.45),
			GroundSpeedKPH:   NewFloat64(56.4),
		},
	},
	{
		name: "bad true track",
		raw:  "$GPVTG,T,45.5,67.5,M,30.45,N,56.40,K*4B",
		err:  "nmea: GPVTG invalid true track: T",
	},
}

func TestVTG(t *testing.T) {
	for _, tt := range vtgtests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				vtg := m.(VTG)
				vtg.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, vtg)
			}
		})
	}
}

var _ = Describe("VTG", func() {
	var (
		parsed VTG
	)
	Describe("Getting data from a $__VTG sentence", func() {
		BeforeEach(func() {
			parsed = VTG{
				TrueTrack:        NewFloat64(TrueDirectionDegrees),
				MagneticTrack:    NewFloat64(MagneticDirectionDegrees),
				GroundSpeedKPH:   NewFloat64(SpeedOverGroundKPH),
				GroundSpeedKnots: NewFloat64(SpeedOverGroundKnots),
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing true track", func() {
			JustBeforeEach(func() {
				parsed.TrueTrack = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetTrueCourseOverGround()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing magnetic track", func() {
			JustBeforeEach(func() {
				parsed.MagneticTrack = Float64{}
			})
			It("should give a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetMagneticCourseOverGround()
				Expect(err).To(HaveOccurred())
			})
			It("should give a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing speed over ground kph", func() {
			JustBeforeEach(func() {
				parsed.GroundSpeedKPH = Float64{}
			})
			It("should give a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing speed over ground knots", func() {
			JustBeforeEach(func() {
				parsed.GroundSpeedKnots = Float64{}
			})
			It("should give a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			It("should give a valid speed over ground", func() {
				Expect(parsed.GetSpeedOverGround()).To(Float64Equal(SpeedOverGroundMPS, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing speed over ground kph and knots", func() {
			JustBeforeEach(func() {
				parsed.GroundSpeedKPH = Float64{}
				parsed.GroundSpeedKnots = Float64{}
			})
			It("should give a valid true course over ground", func() {
				Expect(parsed.GetTrueCourseOverGround()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
			It("should give a valid magnetic course over ground", func() {
				Expect(parsed.GetMagneticCourseOverGround()).To(Float64Equal(MagneticDirectionRadians, 0.00001))
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetSpeedOverGround()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
