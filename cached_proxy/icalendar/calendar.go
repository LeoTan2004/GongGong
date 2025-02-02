package icalendar

import (
	"strings"
)

type IcsCalendar struct {
	events    []Event
	productID string
	timezone  Timezone
}

func (c *IcsCalendar) ToIcs(timezone *Timezone) string {
	result := strings.Builder{}
	result.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\n")
	if c.productID != "" {
		result.WriteString("PRODID:" + c.productID + "\n")
	}

	if timezone == nil {
		if c.timezone == nil {
			timezone = nil
		} else {
			timezone = &c.timezone
		}
	}
	if timezone != nil {
		result.WriteString(c.timezone.ToIcs(nil))
	}
	for _, e := range c.events {
		result.WriteString(e.ToIcs(timezone))
	}
	result.WriteString("END:VCALENDAR\n")
	return result.String()
}

func (c *IcsCalendar) AddEvent(event Event) {
	c.events = append(c.events, event)
}

func (c *IcsCalendar) SetProductID(productID string) {
	c.productID = productID
}

func (c *IcsCalendar) SetTimezone(timezone Timezone) {
	c.timezone = timezone
}
