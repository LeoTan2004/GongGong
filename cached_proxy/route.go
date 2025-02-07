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

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	var creds Credentials
	creds.Username = r.Form.Get("username")
	creds.Password = r.Form.Get("password")
	err = StudentService.SetStudent(creds.Username, creds.Password, true)
	if err == nil {
		token, err := AccountService.Login(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		resp := map[string]string{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   "315360000",
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return nil
	}

	// 检查Token格式（Bearer <token>）
	token := authHeader[len("Bearer "):]
	if token == "" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
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
	w.Header().Set("Content-Type", "application/json")
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
	AccountHandler           = &AccountGetter{TokenService: TokenService{acc: AccountService}}
)

type AccountGetter struct {
	TokenService
}

type IntrospectionResponse struct {
	Active   bool   `json:"active"`             // token 是否有效
	Username string `json:"username,omitempty"` // 用户名
}

func (a *AccountGetter) GetInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	token := r.Form.Get("token")
	account, err := AccountService.GetAccountByToken(token)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	if account == nil {
		return
	} else if account.Status() != account2.Normal {
		resp = map[string]any{
			"active": false,
		}
	} else {
		resp = &IntrospectionResponse{
			Username: account.AccountID(),
			Active:   account.Status() == account2.Normal,
		}
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

type CalendarGetter[V any] struct {
	TokenService
	info            cache.InformationService[V]
	calendarService cache.InformationService[feign.TeachingCalendar]
	convertFunc     func(*V, *feign.TeachingCalendar) icalendar.Calendar
}

var (
	ExamCalendarHandler = CalendarGetter[feign.ExamList]{
		info:            StudentExamService,
		calendarService: CalendarService,
		convertFunc:     ExamsConvertCalendar,
	}
	CoursesCalendarHandler = CalendarGetter[feign.CourseList]{
		info:            StudentCourseService,
		calendarService: CalendarService,
		convertFunc:     CoursesConvertCalendar,
	}
)

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
	if calendar == nil || info == nil {
		http.Error(w, "Data Updating", http.StatusNonAuthoritativeInfo)
		return
	}
	resp := c.convertFunc(info, calendar)
	if resp == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ics := resp.ToIcs(nil)
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	_, err = fmt.Fprint(w, ics)
	if err != nil {
		return
	}

}

const ExamTimeLayout = "2006-01-02 15:04:05"

func ExamsConvertCalendar(exams *feign.ExamList, _ *feign.TeachingCalendar) icalendar.Calendar {
	if exams == nil || exams.Exams == nil {
		return nil
	}
	ical := &icalendar.IcsCalendar{}
	ical.SetProductID(ProdID)
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
		event.SetSummary(fmt.Sprintf("%s %s", ExamSummaryPrefix, exam.Name))
		event.SetLocation(location)
		event.SetDescription(fmt.Sprintf("%s %s %s", exam.Name, exam.Location, ExamDescSuffix))
		event.SetStart(startTime)
		event.SetEnd(endTime)
		for _, a := range DefaultExamAlarms {
			event.AddAlarm(a)
		}
		ical.AddEvent(event)
	}
	return ical
}

func CoursesConvertCalendar(list *feign.CourseList, calendar *feign.TeachingCalendar) icalendar.Calendar {
	if list == nil || list.Courses == nil || calendar == nil {
		return nil
	}
	ical := icalendar.IcsCalendar{}
	ical.SetProductID(ProdID)
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
	summary := fmt.Sprintf("%s %s", CourseSummaryPrefix, course.Name)
	desc := fmt.Sprintf("授课教师：%s  %d节课\\n周次：%s\\n%s", course.Teacher, course.Duration, course.Weeks, CourseDescSummarySuffix)
	location := &icalendar.IcsLocation{}
	location.SetName(course.Classroom)
	event := icalendar.IcsEvent{}
	event.SetSummary(summary)
	event.SetDescription(desc)
	event.SetLocation(location)
	date := calendar.StartTime().AddDate(0, 0, (start-1)*7+feign.Days2Int[course.Day])
	tb := timetable.EventTimes
	startTime := tb[course.StartTime-1].StartTime.AddDate(date.Year(), int(date.Month()), date.Day())
	endTime := tb[course.StartTime+course.Duration-2].EndTime.AddDate(date.Year(), int(date.Month()), date.Day())
	event.SetStart(startTime)
	event.SetEnd(endTime)
	rrule := &icalendar.IcsRepeatRule{}
	rrule.SetFrequency("WEEKLY")
	rrule.SetInterval(1)
	rrule.SetCount(end - start + 1)
	event.SetRepeatRule(rrule)
	for _, a := range DefaultCourseAlarms {
		event.AddAlarm(a)
	}
	return &event
}

func CalPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := w.Write(CalBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
