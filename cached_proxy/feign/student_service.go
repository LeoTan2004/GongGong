package feign

import (
	"cached_proxy/repo"
	"fmt"
	"strings"
)

type StudentService interface {
	// GetStudent 获取学生代理类
	GetStudent(username string) (*Student, error)

	// SetStudent 设置学生账户，如果该账户未通过验证，则返回错误
	SetStudent(username string, password string) error
}

type StudentServiceImpl struct {
	repo   repo.KVRepo[string, *Student]
	client SpiderClient
}

// NewStudentServiceImpl 创建一个新的学生服务实例。client 参数可以为空，此时将使用默认的客户端。
func NewStudentServiceImpl(client *SpiderClient) *StudentServiceImpl {
	if client == nil {
		spiderClient := GetDefaultClient("")
		client = &spiderClient
	}
	return &StudentServiceImpl{client: *client}
}

func (s *StudentServiceImpl) GetStudent(username string) (*Student, error) {
	student, found := s.repo.Get(username)
	if !found {
		return nil, fmt.Errorf("student not found: %s", username)
	}
	return student, nil
}

func (s *StudentServiceImpl) SetStudent(username string, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return fmt.Errorf("invalid username or password")
	}
	student, err := s.client.NewStudent(username, password)
	if err != nil {
		return err
	}
	s.repo.Set(username, &student)
	return nil
}
