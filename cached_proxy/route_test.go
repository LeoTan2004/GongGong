package main

import (
	"cached_proxy/feign"
	"cached_proxy/icalendar"
	"testing"
)

func TestCoursesConvertCalendar(t *testing.T) {
	type args struct {
		list     *feign.CourseList
		calendar *feign.TeachingCalendar
	}
	tests := []struct {
		name  string
		args  args
		judge func(calendar icalendar.Calendar) bool
	}{
		{
			name: "Test CoursesConvertCalendar",
			args: args{
				list: &feign.CourseList{
					Courses: []feign.Course{
						{
							Name:      "Test",
							Teacher:   "Test",
							Classroom: "Test",
							Weeks:     "1-14",
							StartTime: 5,
							Duration:  2,
							Day:       "Monday",
						},
					},
				},
				calendar: &feign.TeachingCalendar{
					Start:  "2025-02-17",
					Weeks:  17,
					TermId: "2024-2025-2",
				},
			},
			judge: func(calendar icalendar.Calendar) bool {
				return calendar != nil && calendar.ToIcs(nil) != ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CoursesConvertCalendar(tt.args.list, tt.args.calendar); !tt.judge(got) {
				t.Errorf("CoursesConvertCalendar() = %v", got)
			}
		})
	}
}
