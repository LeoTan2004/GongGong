package icalendar

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestIcsCalendar_ToIcs(t *testing.T) {
	calendar := &IcsCalendar{}
	calendar.SetProductID("productID")
	result := calendar.ToIcs(nil)
	expected := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:productID\nEND:VCALENDAR\n"
	if result != expected {
		t.Errorf("IcsCalendar.ToIcs() = %v, want %v", result, expected)
	}
}

func TestIcsCalendar_ToIcs2(t *testing.T) {
	calendar := &IcsCalendar{}
	calendar.SetProductID("productID")
	alarms := []Alarm{
		&IcsAlarm{
			action:      "DISPLAY",
			trigger:     -30 * time.Minute,
			description: "display alarm",
		},
		&IcsAlarm{
			action:      "AUDIO",
			trigger:     -40 * time.Minute,
			description: "audio alarm",
		},
	}
	rrule := &IcsRepeatRule{
		frequency: "WEEKLY",
		interval:  1,
		count:     4,
	}
	events := []Event{
		&IcsEvent{
			summary:     "summary",
			description: "description",
			location: &IcsLocation{
				name: "Beijing",
			},
			start:      time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2025, 2, 1, 1, 0, 0, 0, time.UTC),
			alarms:     alarms,
			repeatRule: rrule,
		},
		&IcsEvent{
			summary:     "summary",
			description: "description",
			location: &IcsLocation{
				name: "Beijing",
			},
			start:      time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
			end:        time.Date(2025, 2, 3, 1, 0, 0, 0, time.UTC),
			alarms:     alarms,
			repeatRule: rrule,
		},
	}
	calendar.events = events
	calendar.SetTimezone(&IcsTimezone{
		id:         "Asia/Shanghai",
		name:       "Asia/Shanghai",
		offsetFrom: 8 * time.Hour,
		offsetTo:   8 * time.Hour,
		start:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)})
	result := calendar.ToIcs(nil)
	expectedFilePath := path.Join("testdata", "test_calendar_01.ics")
	bytes, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Errorf("IcsCalendar.ToIcs() = %v, want %v", err, nil)
	}
	expected := string(bytes)
	if result != expected {
		t.Errorf("IcsCalendar.ToIcs() = %v, want %v", result, expected)
	}
}
