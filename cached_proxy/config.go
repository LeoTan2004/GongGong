package main

import (
	"cached_proxy/icalendar"
	"os"
	"time"
)

// SpiderUrl 爬虫地址
var (
	SpiderUrl = os.Getenv("SPIDER_URL")
)

// 日历事件的默认提醒
var (
	// DefaultCourseAlarms 课程事件的默认提醒
	DefaultCourseAlarms = []icalendar.Alarm{
		icalendar.NewIcsAlarm("DISPLAY", 28*time.Minute, "距离上课仅剩28分钟"),
	}
	// DefaultExamAlarms 考试事件的默认提醒
	DefaultExamAlarms = []icalendar.Alarm{
		icalendar.NewIcsAlarm("DISPLAY", 1*time.Hour, "距离考试仅剩1小时"),
		icalendar.NewIcsAlarm("DISPLAY", 24*time.Hour, "距离考试仅剩1天"),
		icalendar.NewIcsAlarm("DISPLAY", 7*24*time.Hour, "距离考试仅剩7天"),
	}
)

// 日历事件的标题和描述的配置
const (
	ExamSummaryPrefix       = "【考试】"
	ExamDescSuffix          = "【拱拱】"
	CourseSummaryPrefix     = "【课程】"
	CourseDescSummarySuffix = "【拱拱】"
	ProdID                  = "-//sky31studio//GongGong//CN"
)

const (
	ApiPort = 8080
)
