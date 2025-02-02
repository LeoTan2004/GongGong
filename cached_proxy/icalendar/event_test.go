package icalendar

import (
	"testing"
	"time"
)

func TestIcsEvent_ToIcs(t *testing.T) {
	type fields struct {
		summary     string
		description string
		location    Location
		start       time.Time
		end         time.Time
		alarms      []Alarm
		repeatRule  RepeatRule
	}
	type args struct {
		timezone *Timezone
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test IcsEvent ToIcs ",
			fields: fields{
				summary:     "summary",
				description: "description",
				start:       time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				end:         time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
			},
			args: args{},
			want: "BEGIN:VEVENT\nSUMMARY:summary\nDESCRIPTION:description\nDTSTART:20210101T000000Z\nDTEND:20210101T010000Z\nUID:summary20210101T000000\nEND:VEVENT\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &IcsEvent{
				summary:     tt.fields.summary,
				description: tt.fields.description,
				location:    tt.fields.location,
				start:       tt.fields.start,
				end:         tt.fields.end,
				alarms:      tt.fields.alarms,
				repeatRule:  tt.fields.repeatRule,
			}
			if got := e.ToIcs(tt.args.timezone); got != tt.want {
				t.Errorf("ToIcs() = %v, want %v", got, tt.want)
			}
		})
	}
}
