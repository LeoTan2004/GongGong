package feign

// SpiderClient 定义了用于与 spider 服务交互的接口。
// 提供了获取教学日历、教室状态、学生信息，以及处理身份验证和其他学生相关操作的方法。
type SpiderClient interface {
	// GetTeachingCalendar 获取当前学期的教学日历。
	// token: 服务的身份验证令牌。
	GetTeachingCalendar(token string) (*TeachingCalendar, error)

	// GetClassroomStatus 获取指定日期的教室考试状态。
	// token: 服务的身份验证令牌。
	// day: 要查询的具体日期（例如，0 表示今天，-1 表示昨天）。
	GetClassroomStatus(token string, day int) (*ClassroomStatusTable, error)

	// GetStudentCourses 获取已认证学生的课程信息。
	// token: 服务的身份验证令牌。
	GetStudentCourses(token string) (*CourseList, error)

	// GetStudentExams 获取已认证学生的考试安排。
	// token: 服务的身份验证令牌。
	GetStudentExams(token string) (*ExamList, error)

	// GetStudentInfo 获取已认证学生的个人信息。
	// token: 服务的身份验证令牌。
	GetStudentInfo(token string) (*StudentInfo, error)

	// Login 使用用户名和密码进行身份验证。
	// username: 用户的用户名。
	// password: 用户的密码。
	Login(username string, password string) (LoginResponse, error)

	// GetStudentScore 获取已认证学生的成绩信息。
	// token: 服务的身份验证令牌。
	// isMajor: 是否仅获取主修相关的成绩。
	GetStudentScore(token string, isMajor bool) (*ScoreBoard, error)

	// GetStudentRank 获取已认证学生的排名信息。
	// token: 服务的身份验证令牌。
	// onlyRequired: 是否仅包括必修课程的排名计算。
	GetStudentRank(token string, onlyRequired bool) (*RankBoard, error)
}
