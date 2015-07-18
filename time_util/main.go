// Package time_util provides a func for getting start of week of given time.
package time_util

import "time"

// StartOfWeek returns time at start of week of t.
func StartOfWeek(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location())
}
