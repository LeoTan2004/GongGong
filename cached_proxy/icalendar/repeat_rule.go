package icalendar

import (
	"fmt"
	"time"
)

type IcsRepeatRule struct {
	frequency string
	interval  int
	count     int
	until     time.Time
}

func (r *IcsRepeatRule) ToIcs(_ *Timezone) string {
	result := "RRULE:FREQ=" + r.frequency
	if r.interval > 1 {
		result += ";INTERVAL=" + string(rune(r.interval))
	}
	if r.count > 0 {
		result += fmt.Sprintf(";COUNT=%d", r.count)
	}
	if !r.until.IsZero() {
		result += ";UNTIL" + TimeToIcs(r.until, nil, "=")
	}
	return result
}

func (r *IcsRepeatRule) SetFrequency(frequency string) {
	r.frequency = frequency
}

func (r *IcsRepeatRule) SetInterval(interval int) {
	r.interval = interval
}

func (r *IcsRepeatRule) SetCount(count int) {
	r.count = count
}

func (r *IcsRepeatRule) SetUntil(until time.Time) {
	r.until = until
}
