package icalendar

import (
	"testing"
	"time"
)

func TestIcsTimezone_ToIcs(t *testing.T) {
	type fields struct {
		id         string
		offsetFrom time.Duration
		offsetTo   time.Duration
		start      time.Time
		name       string
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
			name: "Test IcsTimezone ToIcs",
			fields: fields{
				id:         "id",
				offsetFrom: time.Hour,
				offsetTo:   time.Hour,
				start:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				name:       "name",
			},
			args: args{
				timezone: nil,
			},
			want: "BEGIN:VTIMEZONE\nTZID:id\nBEGIN:STANDARD\nDTSTART:20210101T000000Z\nTZOFFSETFROM:+0100\nTZOFFSETTO:+0100\nTZNAME:name\nEND:STANDARD\nEND:VTIMEZONE\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tz := &IcsTimezone{
				id:         tt.fields.id,
				offsetFrom: tt.fields.offsetFrom,
				offsetTo:   tt.fields.offsetTo,
				start:      tt.fields.start,
				name:       tt.fields.name,
			}
			if got := tz.ToIcs(tt.args.timezone); got != tt.want {
				t.Errorf("ToIcs() = %v, want %v", got, tt.want)
			}
		})
	}
}
