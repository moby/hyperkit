package time_util

import "time"

// StartOfWeek returns time at start of week of t, in local timezone.
func StartOfWeek(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, time.Local)
}
