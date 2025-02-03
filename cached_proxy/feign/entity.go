package feign

import "time"

type CommonResponse[dataType any] struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    dataType `json:"data"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Course struct {
	Name      string `json:"name"`
	Teacher   string `json:"teacher"`
	Classroom string `json:"classroom"`
	Weeks     string `json:"weeks"`
	StartTime int    `json:"start_time"`
	Duration  int    `json:"duration"`
	Day       string `json:"day"`
}

type CourseList struct {
	Courses []Course `json:"courses"`
}

type TeachingCalendar struct {
	Start     string `json:"start"`
	Weeks     int    `json:"weeks"`
	TermId    string `json:"term_id"`
	startTime time.Time
}

func (t *TeachingCalendar) StartTime() time.Time {
	if t.startTime.IsZero() {
		t.startTime, _ = time.Parse("2006-01-02", t.Start)
	}
	return t.startTime
}

type ClassroomStatus struct {
	Name   string   `json:"name"`
	Status []string `json:"status"`
}

type ClassroomStatusTable struct {
	Classrooms map[string][]ClassroomStatus `json:"classrooms"`
	Date       string                       `json:"date"`
}

type Examination struct {
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
	Type      string `json:"type"`
}

type ExamList struct {
	Exams []Examination `json:"exams"`
}

type StudentInfo struct {
	StudentId   string `json:"student_id"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Birthday    string `json:"birthday"`
	Major       string `json:"major"`
	Class       string `json:"class_"`
	EntranceDay string `json:"entrance_day"`
	College     string `json:"college"`
}

type Score struct {
	Name   string `json:"name"`
	Score  string `json:"score"`
	Credit string `json:"credit"`
	Type   string `json:"type"`
	Term   int    `json:"term"`
}

type ScoreBoard struct {
	StudentId         string   `json:"student_id"`
	Name              string   `json:"name"`
	College           string   `json:"college"`
	Major             string   `json:"major"`
	Scores            []Score  `json:"scores"`
	TotalCredit       []string `json:"total_credit"`
	ElectiveCredit    []string `json:"elective_credit"`
	CompulsoryCredit  []string `json:"compulsory_credit"`
	CrossCourseCredit []string `json:"cross_course_credit"`
	AverageScore      string   `json:"average_score"`
	Gpa               string   `json:"gpa"`
	Cet4              string   `json:"cet4"`
	Cet6              string   `json:"cet6"`
}

type Rank struct {
	AverageScore string   `json:"average_score"`
	Gpa          string   `json:"gpa"`
	ClassRank    int      `json:"class_rank"`
	MajorRank    int      `json:"major_rank"`
	Terms        []string `json:"terms"`
}
