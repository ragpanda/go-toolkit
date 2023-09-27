package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDateTzConvert(t *testing.T) {
	now := "2023-06-12 05:01:02"

	cnloc, _ := time.LoadLocation("Asia/Shanghai")
	usEastLoc, _ := time.LoadLocation("America/New_York")
	usWestLoc, _ := time.LoadLocation("America/Los_Angeles")

	cnNow, _ := time.ParseInLocation("2006-01-02 15:04:05", now, cnloc)
	usEastNow, _ := time.ParseInLocation("2006-01-02 15:04:05", now, usEastLoc)
	usWestNow, _ := time.ParseInLocation("2006-01-02 15:04:05", now, usWestLoc)

	dayStart, dayEnd := DayRange(cnNow)
	require.Equal(t, "2023-06-12 00:00:00 +0800 CST", dayStart.String())
	require.Equal(t, "2023-06-13 00:00:00 +0800 CST", dayEnd.String())

	weekStart, weekEnd := WeekRange(cnNow)
	require.Equal(t, "2023-06-12 00:00:00 +0800 CST", weekStart.String())
	require.Equal(t, "2023-06-19 00:00:00 +0800 CST", weekEnd.String())

	dayStart, dayEnd = DayRange(usEastNow)
	require.Equal(t, "2023-06-12 00:00:00 -0400 EDT", dayStart.String())
	require.Equal(t, "2023-06-13 00:00:00 -0400 EDT", dayEnd.String())

	weekStart, weekEnd = WeekRange(usEastNow)
	require.Equal(t, "2023-06-12 00:00:00 -0400 EDT", weekStart.String())
	require.Equal(t, "2023-06-19 00:00:00 -0400 EDT", weekEnd.String())

	dayStart, dayEnd = DayRange(usWestNow)
	require.Equal(t, "2023-06-12 00:00:00 -0700 PDT", dayStart.String())
	require.Equal(t, "2023-06-13 00:00:00 -0700 PDT", dayEnd.String())

	weekStart, weekEnd = WeekRange(usWestNow)
	require.Equal(t, "2023-06-12 00:00:00 -0700 PDT", weekStart.String())
	require.Equal(t, "2023-06-19 00:00:00 -0700 PDT", weekEnd.String())

	now = "2023-06-15 05:01:02"

	cnNow, _ = time.ParseInLocation("2006-01-02 15:04:05", now, cnloc)
	weekStart, weekEnd = WeekRange(cnNow)
	require.Equal(t, "2023-06-12 00:00:00 +0800 CST", weekStart.String())
	require.Equal(t, "2023-06-19 00:00:00 +0800 CST", weekEnd.String())

}
