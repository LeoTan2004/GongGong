package icalendar

import (
	"strings"
	"time"
)

type IcsAlarm struct {
	action      string
	trigger     time.Duration
	description string
}

func NewIcsAlarm(action string, trigger time.Duration, description string) *IcsAlarm {
	return &IcsAlarm{action: action, trigger: trigger, description: description}
}

func (a *IcsAlarm) ToIcs(_ *Timezone) string {
	result := strings.Builder{}
	result.WriteString("BEGIN:VALARM\n")
	if a.action != "" {
		result.WriteString("ACTION:" + a.action + "\n")
	} else {
		result.WriteString("ACTION:DISPLAY\n")
	}
	result.WriteString("TRIGGER:" + DurationToIcs(a.trigger) + "\n")
	result.WriteString("DESCRIPTION:" + a.description + "\n")
	result.WriteString("END:VALARM\n")
	return result.String()
}

func (a *IcsAlarm) SetAction(action string) {
	a.action = action
}

func (a *IcsAlarm) SetTrigger(trigger time.Duration) {
	a.trigger = trigger
}

func (a *IcsAlarm) SetDescription(description string) {
	a.description = description
}
