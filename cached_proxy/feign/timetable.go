package feign

import (
	"math"
	"time"
)

type EventTimes struct {
	StartTime time.Time
	EndTime   time.Time
}

type TimeTable struct {
	EventTimes []EventTimes
}

var (
	zone            = time.FixedZone("Asia/Shanghai", 8*60*60)
	summerStart     = time.Date(0, 5, 1, 0, 0, 0, 0, zone)
	summerEnd       = time.Date(0, 10, 1, 0, 0, 0, 0, zone)
	SummerTimeTable = TimeTable{
		EventTimes: []EventTimes{
			{
				StartTime: time.Date(0, 0, 0, 8, 0, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 8, 45, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 8, 55, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 9, 40, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 10, 10, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 10, 55, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 11, 5, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 11, 50, 0, 0, zone),
			},
			//  下午，这里需要区分冬令时和夏令时
			{
				StartTime: time.Date(0, 0, 0, 14, 30, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 15, 15, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 15, 25, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 16, 10, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 16, 40, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 17, 25, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 17, 35, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 18, 20, 0, 0, zone),
			},
			// 晚上，这里需要区分冬令时和夏令时
			{
				StartTime: time.Date(0, 0, 0, 19, 30, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 20, 15, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 20, 25, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 21, 10, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 21, 20, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 22, 5, 0, 0, zone),
			},
		},
	}
	WinterTimeTable = TimeTable{
		EventTimes: []EventTimes{
			{
				StartTime: time.Date(0, 0, 0, 8, 0, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 8, 45, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 8, 55, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 9, 40, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 10, 10, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 10, 55, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 11, 5, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 11, 50, 0, 0, zone),
			},
			//  下午，这里需要区分冬令时和夏令时
			{
				StartTime: time.Date(0, 0, 0, 14, 0, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 14, 45, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 14, 55, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 15, 40, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 16, 10, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 16, 55, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 17, 5, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 17, 50, 0, 0, zone),
			},
			// 晚上，这里需要区分冬令时和夏令时
			{
				StartTime: time.Date(0, 0, 0, 19, 0, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 19, 45, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 19, 55, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 20, 40, 0, 0, zone),
			}, {
				StartTime: time.Date(0, 0, 0, 20, 50, 0, 0, zone),
				EndTime:   time.Date(0, 0, 0, 21, 35, 0, 0, zone),
			},
		},
	}
)

type TermTimeTable struct {
	SepWeeks     int
	PreTimeTable TimeTable
	SufTimeTable TimeTable
}

func (t *TeachingCalendar) GetTermTimeTable() TermTimeTable {
	sepTime := summerEnd
	termTimeTable := TermTimeTable{
		PreTimeTable: SummerTimeTable,
		SufTimeTable: WinterTimeTable,
	}
	startTime := t.StartTime()
	if startTime.Month() < sepTime.Month() {
		sepTime = summerStart
		termTimeTable.PreTimeTable = WinterTimeTable
		termTimeTable.SufTimeTable = SummerTimeTable
	}
	sepTime = time.Date(startTime.Year(), sepTime.Month(), sepTime.Day(), 0, 0, 0, 0, zone)
	sepWeeks := int(math.Ceil((sepTime.Sub(startTime).Hours())/(24*7))) + 1
	termTimeTable.SepWeeks = sepWeeks
	return termTimeTable
}
