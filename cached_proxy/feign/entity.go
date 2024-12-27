package feign

type CommonResponse[dataType any] struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    dataType `json:"data"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CourseList struct {
	Courses []struct {
		Name      string `json:"name"`
		Teacher   string `json:"teacher"`
		Classroom string `json:"classroom"`
		Weeks     string `json:"weeks"`
		StartTime int    `json:"start_time"`
		Duration  int    `json:"duration"`
		Day       string `json:"day"`
	} `json:"courses"`
}

type TeachingCalendar struct {
	Start  string `json:"start"`
	Weeks  int    `json:"weeks"`
	TermId string `json:"term_id"`
}

type ClassroomStatus struct {
	Name   string   `json:"name"`
	Status []string `json:"status"`
}

type ClassroomStatusTable struct {
	Classrooms map[string][]ClassroomStatus `json:"classrooms"`
	Date       string                       `json:"date"`
}
