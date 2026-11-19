package cli

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/visionik/sogcli/internal/caldav"
)

func captureCalendarOutput(t *testing.T, output func() error) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "output")
	require.NoError(t, err)
	defer file.Close()
	previous := os.Stdout
	os.Stdout = file
	defer func() { os.Stdout = previous }()
	require.NoError(t, output())
	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)
	data, err := io.ReadAll(file)
	require.NoError(t, err)
	return string(data)
}

func TestOutputEventsJSONEscapingAndStatus(t *testing.T) {
	start := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	events := []caldav.Event{
		{UID: "uid\"\\first", Summary: "Review \"budget\"\nSecond line\\path\t\u00e9", Location: "Room \"A\"\\B\nUpstairs", Start: start, End: start.Add(time.Hour), Status: "CANCELLED"},
		{UID: "all-day", Summary: "All day", Start: start, End: start.AddDate(0, 0, 1), AllDay: true},
	}
	output := captureCalendarOutput(t, func() error { return outputEventsJSON(events) })
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	require.Len(t, lines, len(events), "each event must occupy exactly one JSON line")
	for i, line := range lines {
		var got struct {
			UID      string `json:"uid"`
			Summary  string `json:"summary"`
			Location string `json:"location"`
			Start    string `json:"start"`
			End      string `json:"end"`
			AllDay   bool   `json:"all_day"`
			Status   string `json:"status"`
		}
		require.NoError(t, json.Unmarshal([]byte(line), &got))
		require.Equal(t, events[i].UID, got.UID)
		require.Equal(t, events[i].Summary, got.Summary)
		require.Equal(t, events[i].Location, got.Location)
		require.Equal(t, events[i].Start.Format(time.RFC3339), got.Start)
		require.Equal(t, events[i].End.Format(time.RFC3339), got.End)
		require.Equal(t, events[i].AllDay, got.AllDay)
		require.Equal(t, events[i].Status, got.Status)
	}
}

func TestOutputEventsJSONPreservesCancelledStatus(t *testing.T) {
	output := captureCalendarOutput(t, func() error {
		return outputEventsJSON([]caldav.Event{{UID: "cancelled", Summary: "Cancelled meeting", Status: "CANCELLED"}})
	})
	var got map[string]any
	require.NoError(t, json.Unmarshal([]byte(output), &got))
	require.Equal(t, "CANCELLED", got["status"])
	require.Equal(t, false, got["all_day"])
}

func TestCalendarJSONWriteFailure(t *testing.T) {
	for name, output := range map[string]func() error{
		"events":    func() error { return outputEventsJSON([]caldav.Event{{UID: "example"}}) },
		"calendars": func() error { return outputCalendarsJSON([]caldav.Calendar{{Path: "/calendar/"}}) },
	} {
		t.Run(name, func(t *testing.T) {
			file, err := os.CreateTemp(t.TempDir(), "closed-output")
			require.NoError(t, err)
			require.NoError(t, file.Close())
			previous := os.Stdout
			os.Stdout = file
			defer func() { os.Stdout = previous }()
			require.Error(t, output())
		})
	}
}

func TestOutputCalendarsJSONEscaping(t *testing.T) {
	calendar := caldav.Calendar{Path: "/calendars/quoted\"path/", Name: "Team \"Planning\"", Description: "First line\nSecond\\line"}
	output := captureCalendarOutput(t, func() error { return outputCalendarsJSON([]caldav.Calendar{calendar}) })
	require.Equal(t, 1, strings.Count(output, "\n"))
	var got caldav.Calendar
	require.NoError(t, json.Unmarshal([]byte(output), &got))
	require.Equal(t, calendar, got)
}
