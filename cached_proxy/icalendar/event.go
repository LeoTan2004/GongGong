package icalendar

import (
	"strings"
	"time"
)

type IcsEvent struct {
	summary     string
	description string
	location    Location
	start       time.Time
	end         time.Time
	alarms      []Alarm
	repeatRule  RepeatRule
}

func (e *IcsEvent) uid() string {
	return e.summary + e.start.Format("20060102T150405")
}

func (e *IcsEvent) ToIcs(timezone *Timezone) string {
	result := strings.Builder{}
	result.WriteString("BEGIN:VEVENT\n")
	result.WriteString("SUMMARY:" + e.summary + "\n")
	if e.description != "" {
		result.WriteString("DESCRIPTION:" + e.description + "\n")
	}
	if e.location != nil {
		result.WriteString(e.location.ToIcs(timezone) + "\n")
	}
	if !e.start.IsZero() {
		result.WriteString("DTSTART" + TimeToIcs(e.start, timezone, ":") + "\n")
	}
	if !e.end.IsZero() {
		result.WriteString("DTEND" + TimeToIcs(e.end, timezone, ":") + "\n")
	}
	for _, alarm := range e.alarms {
		result.WriteString(alarm.ToIcs(timezone))
	}
	if e.repeatRule != nil {
		result.WriteString(e.repeatRule.ToIcs(timezone) + "\n")
	}
	result.WriteString("UID:" + e.uid() + "\n")
	result.WriteString("END:VEVENT\n")
	return result.String()
}

func (e *IcsEvent) SetSummary(summary string) {
	e.summary = summary
}

func (e *IcsEvent) SetDescription(description string) {
	e.description = description
}

func (e *IcsEvent) SetLocation(location Location) {
	e.location = location
}

func (e *IcsEvent) SetStart(start time.Time) {
	e.start = start
}

func (e *IcsEvent) SetEnd(end time.Time) {
	e.end = end
}

func (e *IcsEvent) AddAlarm(alarm Alarm) {
	e.alarms = append(e.alarms, alarm)
}

func (e *IcsEvent) SetRepeatRule(rule RepeatRule) {
	e.repeatRule = rule
}
