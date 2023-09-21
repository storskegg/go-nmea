package nmea

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/BertoldVdb/go-ais"
	"github.com/BertoldVdb/go-ais/aisnmea"
	"github.com/martinlindhe/unit"
)

const (
	// TypeVDM type for VDM sentences
	TypeVDM = "VDM"

	// TypeVDO type for VDO sentences
	TypeVDO = "VDO"

	rateOfTurnNotAvailable             int16               = -128
	rateOfTurnMaxRightDegreesPerMinute int16               = 127
	rateOfTurnMaxRightRadiansPerSecond float64             = 0.0206
	rateOfTurnMaxLeftDegreesPerMinute  int16               = -rateOfTurnMaxRightDegreesPerMinute
	rateOfTurnMaxLeftRadiansPerSecond  float64             = -rateOfTurnMaxRightRadiansPerSecond
	speedOverGroundNotAvailable        ais.Field10         = 102.3
	latitudeNotAvailable               ais.FieldLatLonFine = 91
	longitudeNotAvailable              ais.FieldLatLonFine = 181
	trueHeadingNotAvailable            uint16              = 511
	cogNotAvailable                    ais.Field10         = 360
)

var navigationStatuses []string = []string{
	"motoring",
	"anchored",
	"not under command",
	"restricted maneuverability",
	"constrained by draft",
	"moored",
	"aground",
	"fishing",
	"sailing",
	"hazardous material high speed",
	"hazardous material wing in ground",
	"reserved for future use",
	"reserved for future use",
	"reserved for future use",
	"ais-sart",
	"default",
}

var vesselTypes []string = []string{
	"Not available (default)",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Reserved for future use",
	"Wing in ground (WIG), All ships of this type",
	"Wing in ground (WIG), Hazardous category A",
	"Wing in ground (WIG), Hazardous category B",
	"Wing in ground (WIG), Hazardous category C",
	"Wing in ground (WIG), Hazardous category D",
	"Wing in ground (WIG), Reserved for future use",
	"Wing in ground (WIG), Reserved for future use",
	"Wing in ground (WIG), Reserved for future use",
	"Wing in ground (WIG), Reserved for future use",
	"Wing in ground (WIG), Reserved for future use",
	"Fishing",
	"Towing",
	"Towing: length exceeds 200m or breadth exceeds 25m",
	"Dredging or underwater ops",
	"Diving ops",
	"Military ops",
	"Sailing",
	"Pleasure Craft",
	"Reserved",
	"Reserved",
	"High speed craft (HSC), All ships of this type",
	"High speed craft (HSC), Hazardous category A",
	"High speed craft (HSC), Hazardous category B",
	"High speed craft (HSC), Hazardous category C",
	"High speed craft (HSC), Hazardous category D",
	"High speed craft (HSC), Reserved for future use",
	"High speed craft (HSC), Reserved for future use",
	"High speed craft (HSC), Reserved for future use",
	"High speed craft (HSC), Reserved for future use",
	"High speed craft (HSC), No additional information",
	"Pilot Vessel",
	"Search and Rescue vessel",
	"Tug",
	"Port Tender",
	"Anti-pollution equipment",
	"Law Enforcement",
	"Spare - Local Vessel",
	"Spare - Local Vessel",
	"Medical Transport",
	"Noncombatant ship according to RR Resolution No. 18",
	"Passenger, All ships of this type",
	"Passenger, Hazardous category A",
	"Passenger, Hazardous category B",
	"Passenger, Hazardous category C",
	"Passenger, Hazardous category D",
	"Passenger, Reserved for future use",
	"Passenger, Reserved for future use",
	"Passenger, Reserved for future use",
	"Passenger, Reserved for future use",
	"Passenger, No additional information",
	"Cargo, All ships of this type",
	"Cargo, Hazardous category A",
	"Cargo, Hazardous category B",
	"Cargo, Hazardous category C",
	"Cargo, Hazardous category D",
	"Cargo, Reserved for future use",
	"Cargo, Reserved for future use",
	"Cargo, Reserved for future use",
	"Cargo, Reserved for future use",
	"Cargo, No additional information",
	"Tanker, All ships of this type",
	"Tanker, Hazardous category A",
	"Tanker, Hazardous category B",
	"Tanker, Hazardous category C",
	"Tanker, Hazardous category D",
	"Tanker, Reserved for future use",
	"Tanker, Reserved for future use",
	"Tanker, Reserved for future use",
	"Tanker, Reserved for future use",
	"Tanker, No additional information",
	"Other Type, All ships of this type",
	"Other Type, Hazardous category A",
	"Other Type, Hazardous category B",
	"Other Type, Hazardous category C",
	"Other Type, Hazardous category D",
	"Other Type, Reserved for future use",
	"Other Type, Reserved for future use",
	"Other Type, Reserved for future use",
	"Other Type, Reserved for future use",
	"Other Type, No additional information",
}

// VDMVDO is a format used to encapsulate generic binary payloads. It is most commonly used
// with AIS data.
// https://gpsd.gitlab.io/gpsd/AIVDM.html
type VDMVDO struct {
	BaseSentence
	NumFragments   Int64
	FragmentNumber Int64
	MessageID      Int64
	Channel        String
	Payload        []byte
	ais.Packet
}

var (
	nmeaCodec *aisnmea.NMEACodec
	aisCodec  *ais.Codec
)

func init() {
	aisCodec = ais.CodecNew(false, false)
	aisCodec.DropSpace = true
	nmeaCodec = aisnmea.NMEACodecNew(aisCodec)
}

// newVDMVDO constructor
func newVDMVDO(s BaseSentence) (VDMVDO, error) {
	p := NewParser(s)
	m := VDMVDO{
		BaseSentence:   s,
		NumFragments:   p.Int64(0, "number of fragments"),
		FragmentNumber: p.Int64(1, "fragment number"),
		MessageID:      p.Int64(2, "sequence number"),
		Channel:        p.String(3, "channel ID"),
		Payload:        p.SixBitASCIIArmour(4, int(p.Int64(5, "number of padding bits").Value), "payload"),
	}
	result, err := nmeaCodec.ParseSentence(s.String())
	if err != nil {
		return m, err
	}
	if result != nil {
		m.Packet = result.Packet
	}
	return m, p.Err()
}

func extractNumber(binaryData []byte, offset int, length int) (uint64, error) {
	if offset < 0 || length < 1 || offset+length >= len(binaryData) {
		return 0, fmt.Errorf("index out of bounds, length of binary data: %d, offset: %d, length: %d", len(binaryData), offset, length)
	}

	var result uint64 = 0

	for _, value := range binaryData[offset : offset+length] {
		result <<= 1
		result |= uint64(value)
	}

	return result, nil
}

func extractString(binaryData []byte, offset int, length int) (string, error) {
	if offset < 0 || length < 1 || offset+length >= len(binaryData) {
		return "", fmt.Errorf("index out of bounds, length of binary data: %d, offset: %d, length: %d", len(binaryData), offset, length)
	}

	if (length)%6 != 0 {
		return "", errors.New("length must be divisible by 6")
	}
	sixBitCharacters := make([]byte, length/6)
	var position int
	for index, value := range binaryData[offset : offset+length] {
		position = index / 6
		sixBitCharacters[position] <<= 1
		sixBitCharacters[position] |= value
	}
	for index, value := range sixBitCharacters {
		if value < 32 {
			sixBitCharacters[index] = value + 64
		}
	}
	return string(sixBitCharacters), nil
}

// GetCallSign retrieves the call sign of the vessel from the sentence
func (s VDMVDO) GetCallSign() (string, error) {
	if staticDataReport, ok := s.Packet.(ais.StaticDataReport); ok && staticDataReport.Valid && staticDataReport.ReportB.Valid {
		return staticDataReport.ReportB.CallSign, nil
	}
	if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		return shipStaticData.CallSign, nil
	}
	return "", fmt.Errorf("value is unavailable")
}

// GetENINumber retrieves the ENI number of the vessel from the sentence
func (s VDMVDO) GetENINumber() (string, error) {
	if binaryBroadcastMessage, ok := s.Packet.(ais.BinaryBroadcastMessage); ok && binaryBroadcastMessage.Valid && binaryBroadcastMessage.ApplicationID.DesignatedAreaCode == 200 && binaryBroadcastMessage.ApplicationID.FunctionIdentifier == 10 {
		eniNumber, err := extractString(binaryBroadcastMessage.BinaryData, 0, 48)
		if err != nil {
			return "", fmt.Errorf("value is unavailable")
		}
		return eniNumber, nil
	}
	return "", fmt.Errorf("value is unavailable")
}

// GetIMONumber retrieves the IMO number of the vessel from the sentence
func (s VDMVDO) GetIMONumber() (string, error) {
	if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		return fmt.Sprintf("%d", shipStaticData.ImoNumber), nil
	}
	return "", fmt.Errorf("value is unavailable")
}

// GetMMSI retrieves the MMSI of the vessel from the sentence
func (s VDMVDO) GetMMSI() (string, error) {
	if s.Packet == nil || s.Packet.GetHeader() == nil {
		return "", fmt.Errorf("value is unavailable")
	}
	return fmt.Sprintf("%d", s.Packet.GetHeader().UserID), nil
}

// GetNavigationStatus retrieves the navigation status from the sentence
func (s VDMVDO) GetNavigationStatus() (string, error) {
	if positionReport, ok := s.Packet.(ais.PositionReport); ok && positionReport.Valid {
		return navigationStatuses[positionReport.NavigationalStatus], nil
	}
	return "", fmt.Errorf("value is unavailable")
}

// GetVesselBeam retrieves the beam of the vessel from the sentence
func (s VDMVDO) GetVesselBeam() (float64, error) {
	if binaryBroadcastMessage, ok := s.Packet.(ais.BinaryBroadcastMessage); ok && binaryBroadcastMessage.Valid && binaryBroadcastMessage.ApplicationID.DesignatedAreaCode == 200 && binaryBroadcastMessage.ApplicationID.FunctionIdentifier == 10 {
		beam, err := extractNumber(binaryBroadcastMessage.BinaryData, 61, 10)
		if err != nil {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Length(beam) * unit.Decimeter).Meters(), nil
	}
	if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		return float64(shipStaticData.Dimension.C + shipStaticData.Dimension.D), nil
	} else if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok {
		return float64(positionReport.Dimension.C + positionReport.Dimension.D), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetVesselLength retrieves the length of the vessel from the sentence
func (s VDMVDO) GetVesselLength() (float64, error) {
	if binaryBroadcastMessage, ok := s.Packet.(ais.BinaryBroadcastMessage); ok && binaryBroadcastMessage.Valid && binaryBroadcastMessage.ApplicationID.DesignatedAreaCode == 200 && binaryBroadcastMessage.ApplicationID.FunctionIdentifier == 10 {
		length, err := extractNumber(binaryBroadcastMessage.BinaryData, 48, 13)
		if err != nil {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Length(length) * unit.Decimeter).Meters(), nil
	}
	if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		return float64(shipStaticData.Dimension.A + shipStaticData.Dimension.B), nil
	} else if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok {
		return float64(positionReport.Dimension.A + positionReport.Dimension.B), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetVesselName retrieves the name of the vessel from the sentence
func (s VDMVDO) GetVesselName() (string, error) {
	if staticDataReport, ok := s.Packet.(ais.StaticDataReport); ok && staticDataReport.Valid && staticDataReport.ReportA.Valid {
		return staticDataReport.ReportA.Name, nil
	}
	if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		return shipStaticData.Name, nil
	}
	if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok && positionReport.Valid {
		return positionReport.Name, nil
	}
	return "", fmt.Errorf("value is unavailable")
}

// GetVesselType retrieves the type of the vessel from the sentence
func (s VDMVDO) GetVesselType() (string, error) {
	vesselTypeIndex := -1
	if staticDataReport, ok := s.Packet.(ais.StaticDataReport); ok && staticDataReport.Valid && staticDataReport.ReportB.Valid {
		vesselTypeIndex = int(staticDataReport.ReportB.ShipType)
	} else if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		vesselTypeIndex = int(shipStaticData.Type)
	} else if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok {
		vesselTypeIndex = int(positionReport.Type)
	}
	if vesselTypeIndex >= 0 && vesselTypeIndex < len(vesselTypes) {
		return vesselTypes[vesselTypeIndex], nil
	}
	return "", fmt.Errorf("value is unavailable")
}

// GetRateOfTurn retrieves the rate of turn from the sentence
func (s VDMVDO) GetRateOfTurn() (float64, error) {
	if positionReport, ok := s.Packet.(ais.PositionReport); ok && positionReport.Valid {
		// https://gpsd.gitlab.io/gpsd/AIVDM.html
		if positionReport.RateOfTurn == rateOfTurnNotAvailable {
			return 0, fmt.Errorf("value is unavailable")
		}
		if positionReport.RateOfTurn == 0 {
			return 0, nil
		}
		if positionReport.RateOfTurn == rateOfTurnMaxLeftDegreesPerMinute {
			return rateOfTurnMaxLeftRadiansPerSecond, nil
		}
		if positionReport.RateOfTurn == rateOfTurnMaxRightDegreesPerMinute {
			return rateOfTurnMaxRightRadiansPerSecond, nil
		}
		aisDecodedROT := math.Pow(float64(positionReport.RateOfTurn)/4.733, 2)
		if positionReport.RateOfTurn < 0 {
			aisDecodedROT = -aisDecodedROT
		}
		return (unit.Angle(aisDecodedROT) * unit.Degree).Radians() / float64(unit.Minute), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetTrueCourseOverGround retrieves the true course over ground from the sentence
func (s VDMVDO) GetTrueCourseOverGround() (float64, error) {
	if positionReport, ok := s.Packet.(ais.PositionReport); ok && positionReport.Valid {
		if positionReport.Cog == cogNotAvailable {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Angle(positionReport.Cog) * unit.Degree).Radians(), nil
	}
	if positionReport, ok := s.Packet.(ais.StandardClassBPositionReport); ok && positionReport.Valid {
		if positionReport.Cog == cogNotAvailable {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Angle(positionReport.Cog) * unit.Degree).Radians(), nil
	}
	if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok && positionReport.Valid {
		if positionReport.Cog == cogNotAvailable {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Angle(positionReport.Cog) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetTrueHeading retrieves the true heading from the sentence
func (s VDMVDO) GetTrueHeading() (float64, error) {
	if positionReport, ok := s.Packet.(ais.PositionReport); ok && positionReport.Valid {
		if positionReport.TrueHeading == trueHeadingNotAvailable {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Angle(positionReport.TrueHeading) * unit.Degree).Radians(), nil
	}
	if positionReport, ok := s.Packet.(ais.StandardClassBPositionReport); ok && positionReport.Valid {
		if positionReport.TrueHeading == trueHeadingNotAvailable {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Angle(positionReport.TrueHeading) * unit.Degree).Radians(), nil
	}
	if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok && positionReport.Valid {
		if positionReport.TrueHeading == trueHeadingNotAvailable {
			return 0, fmt.Errorf("value is unavailable")
		}
		return (unit.Angle(positionReport.TrueHeading) * unit.Degree).Radians(), nil
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetPosition2D retrieves the 2D position from the sentence
func (s VDMVDO) GetPosition2D() (float64, float64, error) {
	if positionReport, ok := s.Packet.(ais.PositionReport); ok && positionReport.Valid {
		if positionReport.Latitude != latitudeNotAvailable && positionReport.Longitude != longitudeNotAvailable {
			return float64(positionReport.Latitude), float64(positionReport.Longitude), nil
		}
	}
	if positionReport, ok := s.Packet.(ais.StandardClassBPositionReport); ok && positionReport.Valid {
		if positionReport.Latitude != latitudeNotAvailable && positionReport.Longitude != longitudeNotAvailable {
			return float64(positionReport.Latitude), float64(positionReport.Longitude), nil
		}
	}
	if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok && positionReport.Valid {
		if positionReport.Latitude != latitudeNotAvailable && positionReport.Longitude != longitudeNotAvailable {
			return float64(positionReport.Latitude), float64(positionReport.Longitude), nil
		}
	}
	return 0, 0, fmt.Errorf("value is unavailable")
}

// GetSpeedOverGround retrieves the speed over ground from the sentence
func (s VDMVDO) GetSpeedOverGround() (float64, error) {
	if positionReport, ok := s.Packet.(ais.PositionReport); ok && positionReport.Valid {
		if positionReport.Sog != speedOverGroundNotAvailable {
			return (unit.Speed(positionReport.Sog) * unit.Knot).MetersPerSecond(), nil
		}
	}
	if positionReport, ok := s.Packet.(ais.StandardClassBPositionReport); ok && positionReport.Valid {
		if positionReport.Sog != speedOverGroundNotAvailable {
			return (unit.Speed(positionReport.Sog) * unit.Knot).MetersPerSecond(), nil
		}
	}
	if positionReport, ok := s.Packet.(ais.ExtendedClassBPositionReport); ok && positionReport.Valid {
		if positionReport.Sog != speedOverGroundNotAvailable {
			return (unit.Speed(positionReport.Sog) * unit.Knot).MetersPerSecond(), nil
		}
	}
	return 0, fmt.Errorf("value is unavailable")
}

// GetDestination retrieves the destination from the sentence
func (s VDMVDO) GetDestination() (string, error) {
	if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		return shipStaticData.Destination, nil
	}
	return "", fmt.Errorf("value is unavailable")
}

// GetETA retrieves the estimated time of arrival from the sentence
func (s VDMVDO) GetETA() (time.Time, error) {
	if shipStaticData, ok := s.Packet.(ais.ShipStaticData); ok && shipStaticData.Valid {
		result := time.Date(
			time.Now().UTC().Year(),
			time.Month(shipStaticData.Eta.Month),
			int(shipStaticData.Eta.Day),
			int(shipStaticData.Eta.Hour),
			int(shipStaticData.Eta.Minute),
			0,
			0,
			time.UTC,
		)
		// The year of the ETA is a guess, if the result is more than a half year ago assume the result should be in the future
		if result.Before(time.Now().UTC().AddDate(0, -6, 0)) {
			result = result.AddDate(1, 0, 0)
		}
		return result, nil
	}
	return time.Unix(0, 0), fmt.Errorf("value is unavailable")
}
