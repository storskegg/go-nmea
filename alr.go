package nmea

import "fmt"

const (
	// TypeALR type for ALR sentences
	TypeALR = "ALR"
	// Active alarm
	ActiveALR = "A"
	// Inactive alarm
	InactiveALR = "V"
	// Active alarm
	AcknowledgedALR = "A"
	// Inactive alarm
	UnacknowledgedALR = "V"
)

// ZDA represents date & time data.
// http://aprs.gids.nl/nmea/#zda
type ALR struct {
	BaseSentence
	Time        Time
	Identifier  String
	Condition   String
	State       String
	Description String
}

// newALR constructor
func newALR(s BaseSentence) (ALR, error) {
	p := NewParser(s)
	p.AssertType(TypeALR)
	return ALR{
		BaseSentence: s,
		Time:         p.Time(0, "time"),
		Identifier:   p.String(1, "identifier"),
		Condition:    p.EnumString(2, "condition", ActiveALR, InactiveALR),
		State:        p.EnumString(3, "state", AcknowledgedALR, UnacknowledgedALR),
		Description:  p.String(4, "description"),
	}, p.Err()
}

func (s ALR) GetIdentifier() (string, error) {
	if !s.Identifier.Valid {
		return "", fmt.Errorf("value is unavailable")
	}
	return s.Identifier.Value, nil
}

func (s ALR) IsActive() (bool, error) {
	return !s.Condition.Valid || s.Condition.Value != InactiveALR, nil
}

func (s ALR) IsUnacknowledged() (bool, error) {
	return !s.State.Valid || s.State.Value != AcknowledgedALR, nil
}

func (s ALR) GetDescription() (string, error) {
	if !s.Description.Valid {
		return "", fmt.Errorf("value is unavailable")
	}
	return s.Description.Value, nil
}
