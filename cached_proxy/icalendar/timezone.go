package icalendar

import (
	"strings"
	"time"
)

type IcsTimezone struct {
	id         string
	offsetFrom time.Duration
	offsetTo   time.Duration
	start      time.Time
	name       string
}

var (
	defaultTimezone = IcsTimezone{
		id:         "Asia/Shanghai",
		offsetFrom: 8,
		offsetTo:   8,
		start:      time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		name:       "Asia/Shanghai",
	}
)

func GetDefaultTimezone() *IcsTimezone {
	return &defaultTimezone
}

func (tz *IcsTimezone) ToIcs(_ *Timezone) string {
	result := strings.Builder{}
	result.WriteString("BEGIN:VTIMEZONE\n")
	result.WriteString("TZID:" + tz.id + "\n")
	result.WriteString("BEGIN:STANDARD\n")
	result.WriteString("DTSTART" + TimeToIcs(tz.start, nil, ":") + "\n")
	result.WriteString("TZOFFSETFROM:" + OffsetToIcs(tz.offsetFrom) + "\n")
	result.WriteString("TZOFFSETTO:" + OffsetToIcs(tz.offsetTo) + "\n")
	result.WriteString("TZNAME:" + tz.name + "\n")
	result.WriteString("END:STANDARD\n")
	result.WriteString("END:VTIMEZONE\n")
	return result.String()
}

func (tz *IcsTimezone) SetID(id string) {
	tz.id = id
}

func (tz *IcsTimezone) SetOffsetFrom(offset time.Duration) {
	tz.offsetFrom = offset
}

func (tz *IcsTimezone) SetOffsetTo(offset time.Duration) {
	tz.offsetTo = offset
}

func (tz *IcsTimezone) SetName(name string) {
	tz.name = name
}

func (tz *IcsTimezone) SetStart(start time.Time) {
	tz.start = start
}

func (tz *IcsTimezone) GetID() string {
	return tz.id
}
