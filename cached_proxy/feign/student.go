package feign

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

type Student interface {
	// GetTeachingCalendar 获取当前学期的教学日历。
	GetTeachingCalendar() (any, error)

	// GetClassroomStatus 获取指定日期的教室考试状态。
	// day: 要查询的具体日期（例如，0 表示今天，-1 表示昨天）。
	GetClassroomStatus(day int) (any, error)

	// GetStudentCourses 获取已认证学生的课程信息。
	GetStudentCourses() (any, error)

	// GetStudentExams 获取已认证学生的考试安排。
	GetStudentExams() (any, error)

	// GetInfo 获取已认证学生的个人信息。
	GetInfo() (any, error)

	// GetStudentScore 获取已认证学生的成绩信息。
	// isMajor: 是否仅获取主修相关的成绩。
	GetStudentScore(isMajor bool) (any, error)

	// GetStudentRank 获取已认证学生的排名信息。
	// onlyRequired: 是否仅包括必修课程的排名计算。
	GetStudentRank(onlyRequired bool) (any, error)
}

type StudentImpl struct {
	spider       SpiderClient
	username     string
	password     string
	token        string // 静态token，用于对外部使用，一般不会变化
	dynamicToken string // 动态token，可能遇到失效的问题，会重新申请
}

// GenerateUUID 生成一个UUID
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		return "", err
	}
	// Set the version (4) and variant bits as per RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10xxxxxx

	// Format as a UUID string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16]), nil
}

func NewStudentImpl(username string, password string, client SpiderClient) (*StudentImpl, error) {
	s := StudentImpl{username: username, password: password, spider: client}
	_, err := s.refreshDynamicToken(3)
	if err != nil {
		return nil, err
	}
	s.token, err = GenerateUUID()
	return &s, err
}

func (s *StudentImpl) Username() string {
	return s.username
}

func (s *StudentImpl) AccountId() string {
	return s.username
}

func (s *StudentImpl) Token() string {
	return s.token
}

// refreshDynamicToken 刷新动态token
func (s *StudentImpl) refreshDynamicToken(maxRetryTime int) (string, error) {
	var finalError error

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

func (s *StudentImpl) doGetter(function func(token string) (any, error)) (any, error) {
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

func (s *StudentImpl) GetTeachingCalendar() (any, error) {
	return s.doGetter(s.spider.GetTeachingCalendar)
}

func (s *StudentImpl) GetClassroomStatus(day int) (any, error) {
	return s.doGetter(func(token string) (any, error) {
		return s.spider.GetClassroomStatus(token, day)
	})
}

func (s *StudentImpl) GetStudentCourses() (any, error) {
	return s.doGetter(s.spider.GetStudentCourses)
}

func (s *StudentImpl) GetStudentExams() (any, error) {
	return s.doGetter(s.spider.GetStudentExams)
}

func (s *StudentImpl) GetInfo() (any, error) {
	return s.doGetter(s.spider.GetStudentInfo)
}

func (s *StudentImpl) GetStudentScore(isMajor bool) (any, error) {
	return s.doGetter(func(token string) (any, error) {
		return s.spider.GetStudentScore(token, isMajor)
	})
}

func (s *StudentImpl) GetStudentRank(onlyRequired bool) (any, error) {
	return s.doGetter(func(token string) (any, error) {
		return s.spider.GetStudentRank(token, onlyRequired)
	})
}
