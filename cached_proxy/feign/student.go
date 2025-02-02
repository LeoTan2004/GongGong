package feign

import (
	"fmt"
	"strings"
	"sync"
)

type Student interface {
	// Username 获取学生的用户名。
	Username() string

	// GetTeachingCalendar 获取当前学期的教学日历。
	GetTeachingCalendar() (*TeachingCalendar, error)

	// GetClassroomStatus 获取指定日期的教室考试状态。
	// day: 要查询的具体日期（例如，0 表示今天，-1 表示昨天）。
	GetClassroomStatus(day int) (*ClassroomStatusTable, error)

	// GetStudentCourses 获取已认证学生的课程信息。
	GetStudentCourses() (*CourseList, error)

	// GetStudentExams 获取已认证学生的考试安排。
	GetStudentExams() (*ExamList, error)

	// GetInfo 获取已认证学生的个人信息。
	GetInfo() (*StudentInfo, error)

	// GetStudentScore 获取已认证学生的成绩信息。
	// isMajor: 是否仅获取主修相关的成绩。
	GetStudentScore(isMajor bool) (*ScoreBoard, error)

	// GetStudentRank 获取已认证学生的排名信息。
	// onlyRequired: 是否仅包括必修课程的排名计算。
	GetStudentRank(onlyRequired bool) (*Rank, error)
}

type StudentImpl struct {
	spider       SpiderClient
	username     string
	password     string
	dynamicToken string // 动态token，可能遇到失效的问题，会重新申请
	version      int    // 版本号，用于防止并发登陆问题
	mu           sync.Mutex
}

func NewStudentImpl(username string, password string, client SpiderClient) (*StudentImpl, error) {
	s := StudentImpl{username: username, password: password, spider: client}
	_, err := s.refreshDynamicToken(3)
	if err != nil {
		return nil, err
	}
	return &s, err
}

func (s *StudentImpl) Username() string {
	return s.username
}

// refreshDynamicToken 刷新动态token
func (s *StudentImpl) refreshDynamicToken(maxRetryTime int) (string, error) {
	var finalError error
	version := s.version
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
		s.version++
	}()
	if version != s.version {
		return s.dynamicToken, nil
	}
	for i := 0; i < maxRetryTime; i++ {
		response, err := s.spider.Login(s.username, s.password)
		if err != nil {
			finalError = err
			errorString := err.Error()
			// 如果是账号密码错误，那么重试
			if strings.Contains(errorString, "unauthorized") {
				return "", err
			}
			// 否则重试
			continue
		}
		s.dynamicToken = response.Token
		return s.dynamicToken, nil
	}
	return "", finalError
}

func doGetter[V any](s *StudentImpl, function func(token string) (*V, error)) (*V, error) {
	var finalErr error
	maxRetryTimes := 3
	// 如果token为空，那么刷新token
	if s.dynamicToken == "" {
		_, err := s.refreshDynamicToken(3)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < maxRetryTimes; i++ {
		token := s.dynamicToken
		data, err := function(token)
		if err != nil {
			finalErr = err
			errorString := err.Error()
			switch errorString {
			case "unauthorized":
				// 如果是token失效，那么重试登陆
				_, err := s.refreshDynamicToken(3)
				if err != nil {
					return nil, err
				}
				continue
			case "service unavailable":
				// 如果是服务不可用，那么重试
				continue
			default:
				// 否则返回错误
				return nil, err
			}
		} else {
			return data, nil
		}
	}
	// 如果重试次数超过限制，那么返回错误
	return nil, fmt.Errorf("exceeded retry attempts: %v", finalErr)
}

func (s *StudentImpl) GetTeachingCalendar() (*TeachingCalendar, error) {
	return doGetter(s, s.spider.GetTeachingCalendar)
}

func (s *StudentImpl) GetClassroomStatus(day int) (*ClassroomStatusTable, error) {
	return doGetter(s, func(token string) (*ClassroomStatusTable, error) {
		return s.spider.GetClassroomStatus(token, day)
	})
}

func (s *StudentImpl) GetStudentCourses() (*CourseList, error) {
	return doGetter(s, s.spider.GetStudentCourses)
}

func (s *StudentImpl) GetStudentExams() (*ExamList, error) {
	return doGetter(s, s.spider.GetStudentExams)
}

func (s *StudentImpl) GetInfo() (*StudentInfo, error) {
	return doGetter(s, s.spider.GetStudentInfo)
}

func (s *StudentImpl) GetStudentScore(isMajor bool) (*ScoreBoard, error) {
	return doGetter(s, func(token string) (*ScoreBoard, error) {
		return s.spider.GetStudentScore(token, isMajor)
	})
}

func (s *StudentImpl) GetStudentRank(onlyRequired bool) (*Rank, error) {
	return doGetter(s, func(token string) (*Rank, error) {
		return s.spider.GetStudentRank(token, onlyRequired)
	})
}
