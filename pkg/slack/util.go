package slack

import (
	"fmt"
	"strconv"
	"strings"
)

// Timestamp returns the seconds part of a Slack timestamp
// Example: "1533028651.000211" -> 1533028651
func Timestamp(ts string) int64 {
	s := strings.Split(ts, ".")
	i, _ := strconv.ParseInt(s[0], 10, 64)
	return i
}

// TimestampNano returns a Slack timestamp as nanoseconds
func TimestampNano(ts string) int64 {
	s := strings.Split(ts, ".")
	i, _ := strconv.ParseInt(s[0], 10, 64)
	j, _ := strconv.ParseInt(s[1], 10, 64)

	return (i * 1000000) + j
}

// TimestampNanoString converts a Slack TS in nanoseconds into Slack's string representation
func TimestampNanoString(ts int64) string {
	_p1 := ts / 1000000
	_p2 := ts % 1000000
	return fmt.Sprintf("%d.%06d", _p1, _p2)
}
