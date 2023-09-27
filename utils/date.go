package utils

import "time"

const Day = 24 * time.Hour

// TruncateDay returns the start of the day of the given time.
func TruncateDay(t time.Time) time.Time {
	year, month, day := t.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return start
}

// TruncateWeek returns the start of the week of the given time.
func TruncateWeek(t time.Time) time.Time {
	t = TruncateDay(t)
	dayStart := InLocWithSameLiteralTime(t, time.UTC)
	weekStart := dayStart.Truncate(7 * Day)
	weekStart = InLocWithSameLiteralTime(weekStart, t.Location())
	return TruncateDay(weekStart)
}

// DayRange returns the start and end of the day of the given time.
func DayRange(t time.Time) (time.Time, time.Time) {
	start := TruncateDay(t)
	return start, start.Add(24 * time.Hour)
}

// WeekRange returns the start and end of the week of the given time.
func WeekRange(t time.Time) (time.Time, time.Time) {
	weekStart := TruncateWeek(t)
	return weekStart, weekStart.Add(7 * Day)
}

// InLocWithSameLiteralTime returns a time with the same literal time in the given location.
// For example, if t is 2019-01-01 00:00:00 +0800 CST, and loc is UTC,
// then the returned time will be 2019-01-01 00:00:00 +0000 UTC.
func InLocWithSameLiteralTime(t time.Time, loc *time.Location) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(year, month, day, hour, min, sec, t.Nanosecond(), loc)
}

// InLocWithSameTrueTime returns a time with the same true time in the given location.
// For example, if t is 2019-01-01 00:00:00 +0800 CST, and loc is UTC,
// then the returned time will be 2019-01-01 08:00:00 +0000 UTC.
func InLocWithSameTrueTime(t time.Time, loc *time.Location) time.Time {
	return t.In(loc)
}
