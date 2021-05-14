package nmea

import (
	"fmt"
	"strings"
)

// TagBlock type
type TagBlock struct {
	Valid         bool
	InvalidReason string
	Time          Int64  // TypeUnixTime unix timestamp (unit is likely to be s, but might be ms, YMMV), parameter: -c
	RelativeTime  Int64  // TypeRelativeTime relative time, parameter: -r
	Destination   String // TypeDestinationID destination identification 15 char max, parameter: -d
	Grouping      String // TypeGrouping sentence grouping, parameter: -g
	LineCount     Int64  // TypeLineCount line count, parameter: -n
	Source        String // TypeSourceID source identification 15 char max, parameter: -s
	Text          String // TypeTextString valid character string, parameter -t
}

// NewInvalidTagblock creates an invalid TagBlock
func NewInvalidTagblock(reason string) TagBlock {
	return TagBlock{
		InvalidReason: reason,
	}
}

// NewTagblock creates a valid TagBlock
func NewTagblock() TagBlock {
	return TagBlock{
		Valid:        true,
		Time:         NewInvalidInt64("not specified"),
		RelativeTime: NewInvalidInt64("not specified"),
		Destination:  NewInvalidString("not specified"),
		Grouping:     NewInvalidString("not specified"),
		LineCount:    NewInvalidInt64("not specified"),
		Source:       NewInvalidString("not specified"),
		Text:         NewInvalidString("not specified"),
	}
}

// parseTagBlock adds support for tagblocks
// https://gpsd.gitlab.io/gpsd/AIVDM.html#_nmea_tag_blocks
func parseTagBlock(tags string) TagBlock {
	sumSepIndex := strings.Index(tags, ChecksumSep)
	if sumSepIndex == -1 {
		return NewInvalidTagblock("nmea: tagblock does not contain checksum separator")
	}

	var (
		fieldsRaw   = tags[0:sumSepIndex]
		checksumRaw = strings.ToUpper(tags[sumSepIndex+1:])
		checksum    = Checksum(fieldsRaw)
		tagBlock    = NewTagblock()
	)

	// Validate the checksum
	if checksum != checksumRaw {
		return NewInvalidTagblock(fmt.Sprintf("nmea: tagblock checksum mismatch [%s != %s]", checksum, checksumRaw))
	}

	items := strings.Split(tags[:sumSepIndex], ",")
	for _, item := range items {
		parts := strings.SplitN(item, ":", 2)
		if len(parts) != 2 {
			return NewInvalidTagblock(fmt.Sprintf("nmea: tagblock field is malformed (should be <key>:<value>) [%s]", item))
		}
		key, value := parts[0], parts[1]
		switch key {
		case "c": // UNIX timestamp
			tagBlock.Time = ParseInt64(value)
		case "d": // Destination ID
			tagBlock.Destination = NewString(value)
		case "g": // Grouping
			tagBlock.Grouping = NewString(value)
		case "n": // Line count
			tagBlock.LineCount = ParseInt64(value)
		case "r": // Relative time
			tagBlock.RelativeTime = ParseInt64(value)
		case "s": // Source ID
			tagBlock.Source = NewString(value)
		case "t": // Text string
			tagBlock.Text = NewString(value)
		}
	}
	return tagBlock
}
