package main

import (
	account2 "cached_proxy/account"
	"cached_proxy/cache"
	"cached_proxy/feign"
	"cached_proxy/icalendar"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Whoami(w http.ResponseWriter, r *http.Request) {
	log.Printf("Whoami %s\n", r.RequestURI)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.Header.Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	account, err := AccountService.GetAccountByToken(token)
	if err != nil {
		if err.Error() == "account not found" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	if account == nil || account.Status() != account2.Normal {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	resp := feign.CommonResponse[map[string]string]{
		Code:    1,
		Message: "success",
		Data: map[string]string{
			"username": account.AccountID(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = StudentService.SetStudent(creds.Username, creds.Password, true)
	if err == nil {
		token, err := AccountService.Login(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		resp := feign.CommonResponse[map[string]string]{
			Code:    1,
			Message: "success",
			Data: map[string]string{
				"token": token,
			},
		}
		// 返回 token
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}
	if err.Error() == "unauthorized" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type TokenService struct {
	acc account2.Service
}

func (t *TokenService) checkToken(w http.ResponseWriter, r *http.Request) account2.Account {
	token := r.Header.Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	account, err := AccountService.GetAccountByToken(token)
	if err != nil {
		if err.Error() == "account not found" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return nil
	}
	if account == nil || account.Status() != account2.Normal {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	return account
}

type InfoGetter[V any] struct {
	TokenService
	info cache.InformationService[V]
}

func (c *InfoGetter[V]) GetInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetInfo %s\n", r.RequestURI)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	account := c.checkToken(w, r)
	if account == nil {
		return
	}
	info, err := c.info.GetInfo(account.AccountID())
	if err != nil {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
	}
	resp := feign.CommonResponse[any]{
		Code:    1,
		Message: "success",
		Data:    info,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

var (
	CalendarHandler          = &InfoGetter[feign.TeachingCalendar]{info: CalendarService}
	TodayClassroomHandler    = &InfoGetter[feign.ClassroomStatusTable]{info: TodayClassroomService}
	TomorrowClassroomHandler = &InfoGetter[feign.ClassroomStatusTable]{info: TomorrowClassroomService}
	InfoHandler              = &InfoGetter[feign.StudentInfo]{info: StudentInfoService}
	MajorScoreHandler        = &InfoGetter[feign.ScoreBoard]{info: StudentMajorScoreService}
	MinorScoreHandler        = &InfoGetter[feign.ScoreBoard]{info: StudentMinorScoreService}
	TotalRankHandler         = &InfoGetter[feign.Rank]{info: StudentTotalRankService}
	RequiredRankHandler      = &InfoGetter[feign.Rank]{info: StudentRequiredRankService}
	ExamHandler              = &InfoGetter[feign.ExamList]{info: StudentExamService}
	CourseHandler            = &InfoGetter[feign.CourseList]{info: StudentCourseService}
)

type CalendarGetter[V any] struct {
	TokenService
	info            cache.InformationService[V]
	calendarService cache.InformationService[feign.TeachingCalendar]
	convertFunc     func(*V, *feign.TeachingCalendar) icalendar.Calendar
}

func (c *CalendarGetter[V]) GetInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetInfo %s\n", r.RequestURI)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	account := c.checkToken(w, r)
	if account == nil {
		return
	}
	info, err := c.info.GetInfo(account.AccountID())
	if err != nil {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
	}
	calendar, err := c.calendarService.GetInfo(account.AccountID())
	if err != nil {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
	}
	resp := c.convertFunc(info, calendar)
	if resp == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ics := resp.ToIcs(nil)
	w.Header().Set("Content-Type", "text/calendar")
	_, err = w.Write([]byte(ics))
	if err != nil {
		return
	}

}

var (
	CourseAlarm = icalendar.NewIcsAlarm("DISPLAY", 25*time.Minute, "距离上课仅剩25分钟")
)

const ExamTimeLayout = "2006-01-02 15:04:05"

func ExamsConvertCalendar(exams *feign.ExamList, calendar *feign.TeachingCalendar) icalendar.Calendar {
	if exams == nil || exams.Exams == nil {
		return nil
	}
	ical := &icalendar.IcsCalendar{}
	ical.SetTimezone(icalendar.GetDefaultTimezone())
	for _, exam := range exams.Exams {
		if exam.StartTime == "" {
			continue
		}
		event := &icalendar.IcsEvent{}
		var startTime, endTime time.Time
		var err error
		if startTime, err = time.Parse(ExamTimeLayout, exam.StartTime); err != nil {
			continue
		}
		if endTime, err = time.Parse(ExamTimeLayout, exam.EndTime); err != nil {
			endTime = startTime
		}
		location := &icalendar.IcsLocation{}
		location.SetName(exam.Location)
		event.SetSummary(fmt.Sprintf("【考试】%s", exam.Name))
		event.SetLocation(location)
		event.SetDescription(fmt.Sprintf("【%s】%s", exam.Name, exam.Location))
		event.SetStart(startTime)
		event.SetEnd(endTime)
		ical.AddEvent(event)
	}
	return ical
}

func CoursesConvertCalendar(list *feign.CourseList, calendar *feign.TeachingCalendar) icalendar.Calendar {
	if list == nil || list.Courses == nil || calendar == nil {
		return nil
	}
	ical := icalendar.IcsCalendar{}
	ical.SetTimezone(icalendar.GetDefaultTimezone())
	timetable := calendar.GetTermTimeTable()
	for _, course := range list.Courses {
		if course.Weeks == "" || course.Day == "" || course.StartTime == 0 || course.Duration == 0 {
			continue
		}
		weeks := strings.Split(course.Weeks, ",")
		for _, week := range weeks {
			week = strings.TrimSpace(week)
			if week == "" {
				continue
			}
			w := strings.Split(week, "-")
			var start, end int
			var err error
			if start, err = strconv.Atoi(w[0]); err != nil {
				log.Printf("failed to parse week: %s", week)
				continue
			}
			if len(w) < 2 {
				end = start
			} else if end, err = strconv.Atoi(w[1]); err != nil {
				log.Printf("failed to parse week: %s", week)
				continue
			}

			//	If the time was cross the sep week, separate the event into two parts
			//  e.g. sep = 11  start = 10 end = 11
			//  e.g. sep = 11  start = 10 end = 12
			//  n.e.g. sep = 11  start = 11 end = 12
			if start < timetable.SepWeeks && end >= timetable.SepWeeks && course.StartTime+course.Duration-1 > 4 {
				event := convertCourseToEvent(course, calendar, start, timetable.SepWeeks-1, timetable.PreTimeTable)
				ical.AddEvent(event)
				start = timetable.SepWeeks
			}
			if end >= timetable.SepWeeks {
				event := convertCourseToEvent(course, calendar, timetable.SepWeeks, end, timetable.SufTimeTable)
				ical.AddEvent(event)
			} else {
				event := convertCourseToEvent(course, calendar, start, end, timetable.PreTimeTable)
				ical.AddEvent(event)
			}
		}
	}
	return &ical
}

func convertCourseToEvent(course feign.Course, calendar *feign.TeachingCalendar, start int, end int, timetable feign.TimeTable) *icalendar.IcsEvent {
	summary := fmt.Sprintf("【课程】%s", course.Name)
	desc := fmt.Sprintf("【%s】%d节课", course.Teacher, course.Duration)
	location := &icalendar.IcsLocation{}
	location.SetName(course.Classroom)
	event := icalendar.IcsEvent{}
	event.SetSummary(summary)
	event.SetDescription(desc)
	event.SetLocation(location)
	date := calendar.StartTime().AddDate(0, 0, (start-1)*7+feign.Days2Int[course.Day])
	tb := timetable.Times
	startTime := tb[course.StartTime-1].StartTime.AddDate(date.Year(), int(date.Month()), date.Day())
	endTime := tb[course.StartTime+course.Duration-2].EndTime.AddDate(date.Year(), int(date.Month()), date.Day())
	event.SetStart(startTime)
	event.SetEnd(endTime)
	rrule := &icalendar.IcsRepeatRule{}
	rrule.SetFrequency("WEEKLY")
	rrule.SetInterval(1)
	rrule.SetCount(end - start + 1)
	event.SetRepeatRule(rrule)
	event.AddAlarm(CourseAlarm)
	return &event
}
