package icalendar

import (
	"testing"
	"time"
)

func TestIcsAlarm_ToIcs(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		trigger  time.Duration
		desc     string
		expected string
	}{
		{
			name:     "TestIcsAlarm_ToIcs_WithDefaultAction",
			action:   "",
			trigger:  0,
			desc:     "Test",
			expected: "BEGIN:VALARM\nACTION:DISPLAY\nTRIGGER:PT0S\nDESCRIPTION:Test\nEND:VALARM\n",
		},
		{
			name:     "TestIcsAlarm_ToIcs_WithCustomAction",
			action:   "AUDIO",
			trigger:  0,
			desc:     "Test",
			expected: "BEGIN:VALARM\nACTION:AUDIO\nTRIGGER:PT0S\nDESCRIPTION:Test\nEND:VALARM\n",
		},
		{
			name:     "TestIcsAlarm_ToIcs_WithCustomTrigger",
			action:   "DISPLAY",
			trigger:  1*time.Hour + 20*time.Second,
			desc:     "Test",
			expected: "BEGIN:VALARM\nACTION:DISPLAY\nTRIGGER:PT1H20S\nDESCRIPTION:Test\nEND:VALARM\n",
		},
		{
			name:     "TestIcsAlarm_ToIcs_WithCustomActionAndTrigger",
			action:   "DISPLAY",
			trigger:  25*time.Hour + 20*time.Second,
			desc:     "Test",
			expected: "BEGIN:VALARM\nACTION:DISPLAY\nTRIGGER:P1DT1H20S\nDESCRIPTION:Test\nEND:VALARM\n",
		},
		{
			name:     "TestIcsAlarm_ToIcs_WithCustomDescription",
			action:   "DISPLAY",
			trigger:  0,
			desc:     "Custom",
			expected: "BEGIN:VALARM\nACTION:DISPLAY\nTRIGGER:PT0S\nDESCRIPTION:Custom\nEND:VALARM\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := IcsAlarm{}
			a.SetAction(tt.action)
			a.SetTrigger(tt.trigger)
			a.SetDescription(tt.desc)
			if got := a.ToIcs(nil); got != tt.expected {
				t.Errorf("IcsAlarm.ToIcs() = %v, want %v", got, tt.expected)
			}
		})
	}
}
