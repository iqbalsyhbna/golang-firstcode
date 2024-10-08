// internal/helpers/time_helpers.go
package helpers

import (
	"fmt"
	"time"
)

const (
	TimestampFormat   = "2006-01-02 15:04:05"
	DateFormat        = "2006-01-02"
	TimestampTZFormat = "2006-01-02 15:04:05 -0700"
)

func ParseTimestamp(b []uint8) (time.Time, error) {
	return time.Parse(TimestampFormat, string(b))
}

func ParseDate(b []uint8) (time.Time, error) {
	return time.Parse(DateFormat, string(b))
}

func ParseTimestampTZ(b []uint8) (time.Time, error) {
	return time.Parse(TimestampTZFormat, string(b))
}

// ParseAnyTimestamp tries different formats
func ParseAnyTimestamp(b []uint8) (time.Time, error) {
	formats := []string{TimestampFormat, DateFormat, TimestampTZFormat}
	str := string(b)

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", str)
}
