package feign

import (
	"reflect"
	"testing"
	"time"
)

func TestTeachingCalendar_GetTermTimeTable(t1 *testing.T) {
	type fields struct {
		Start     string
		Weeks     int
		TermId    string
		startTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   TermTimeTable
	}{
		{
			name: "Test Term 1",
			fields: fields{
				Start:  "2025-02-17",
				Weeks:  18,
				TermId: "2024-2025-2",
			},
			want: TermTimeTable{
				SepWeeks:     12,
				PreTimeTable: WinterTimeTable,
				SufTimeTable: SummerTimeTable,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TeachingCalendar{
				Start:     tt.fields.Start,
				Weeks:     tt.fields.Weeks,
				TermId:    tt.fields.TermId,
				startTime: tt.fields.startTime,
			}
			if got := t.GetTermTimeTable(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetTermTimeTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
