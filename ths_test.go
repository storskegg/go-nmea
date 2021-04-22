package nmea_test

import (
	"testing"

	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var thstests = []struct {
	name string
	raw  string
	err  string
	msg  THS
}{
	{
		name: "good sentence AutonomousTHS",
		raw:  "$INTHS,123.456,A*20",
		msg: THS{
			Heading: NewFloat64(123.456),
			Status:  AutonomousTHS,
		},
	},
	{
		name: "good sentence EstimatedTHS",
		raw:  "$INTHS,123.456,E*24",
		msg: THS{
			Heading: NewFloat64(123.456),
			Status:  EstimatedTHS,
		},
	},
	{
		name: "good sentence ManualTHS",
		raw:  "$INTHS,123.456,M*2C",
		msg: THS{
			Heading: NewFloat64(123.456),
			Status:  ManualTHS,
		},
	},
	{
		name: "good sentence SimulatorTHS",
		raw:  "$INTHS,123.456,S*32",
		msg: THS{
			Heading: NewFloat64(123.456),
			Status:  SimulatorTHS,
		},
	},
	{
		name: "good sentence InvalidTHS",
		raw:  "$INTHS,,V*1E",
		msg: THS{
			Heading: Float64{},
			Status:  InvalidTHS,
		},
	},
	{
		name: "invalid Status",
		raw:  "$INTHS,123.456,B*23",
		err:  "nmea: INTHS invalid status: B",
	},
	{
		name: "invalid Heading",
		raw:  "$INTHS,XXX,A*51",
		err:  "nmea: INTHS invalid heading: XXX",
	},
}

func TestTHS(t *testing.T) {
	for _, tt := range thstests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				ths := m.(THS)
				ths.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, ths)
			}
		})
	}
}

var _ = Describe("THS", func() {
	var (
		parsed THS
	)
	Describe("Getting data from a $__THS sentence", func() {
		BeforeEach(func() {
			parsed = THS{
				Heading: NewFloat64(TrueDirectionDegrees),
				Status:  SimulatorTHS,
			}
		})
		Context("When having a parsed sentence", func() {
			It("should give a valid true heading", func() {
				Expect(parsed.GetTrueHeading()).To(Float64Equal(TrueDirectionRadians, 0.00001))
			})
		})
		Context("When having a parsed sentence with missing heading", func() {
			JustBeforeEach(func() {
				parsed.Heading = Float64{}
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When having a parsed sentence with status flag set to invalid", func() {
			JustBeforeEach(func() {
				parsed.Status = InvalidTHS
			})
			Specify("an error is returned", func() {
				_, err := parsed.GetTrueHeading()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
