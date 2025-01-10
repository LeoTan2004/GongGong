package service

import (
	"cached_proxy/feign"
	"cached_proxy/repo"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type SpiderService interface {
	// GetStudent 获取学生代理类
	GetStudent(username string) (*feign.Student, error)

	// SetStudent 设置学生账户
	SetStudent(username string, password string) error
}

var baseUrl = os.Getenv("SPIDER_BASE_URL") // 爬虫服务的基础URL

var service SpiderService

// GetSpiderService singleton mode
func GetSpiderService() SpiderService {
	if service == nil {
		service = &SpiderServiceImpl{
			repo:   repo.NewMemRepo[string, *feign.Student](),
			client: feign.NewSpiderClientImpl(baseUrl, http.Client{}),
		}
	}
	return service
}

type SpiderServiceImpl struct {
	repo   repo.KVRepo[string, *feign.Student]
	client feign.SpiderClient
}

func (s *SpiderServiceImpl) GetStudent(username string) (*feign.Student, error) {
	student, found := s.repo.Get(username)
	if !found {
		return nil, fmt.Errorf("student not found: %s", username)
	}
	return student, nil
}

func (s *SpiderServiceImpl) SetStudent(username string, password string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return fmt.Errorf("username and password cannot be empty")
	}
	student, err := s.client.NewStudent(username, password)
	if err != nil {
		return err
	}
	s.repo.Set(username, &student)
	return nil
}
