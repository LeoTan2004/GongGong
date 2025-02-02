package icalendar

import (
	"testing"
	"time"
)

func TestIcsRepeatRule_ToIcs(t *testing.T) {
	type fields struct {
		frequency string
		interval  int
		count     int
		until     time.Time
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
			name: "Test IcsRepeatRule ToIcs",
			fields: fields{
				frequency: "DAILY",
				interval:  1,
				count:     0,
				until:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				timezone: nil,
			},
			want: "RRULE:FREQ=DAILY;UNTIL=20210101T000000Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IcsRepeatRule{
				frequency: tt.fields.frequency,
				interval:  tt.fields.interval,
				count:     tt.fields.count,
				until:     tt.fields.until,
			}
			if got := r.ToIcs(tt.args.timezone); got != tt.want {
				t.Errorf("ToIcs() = %v, want %v", got, tt.want)
			}
		})
	}
}
