package main

import (
	"cached_proxy/account"
	"cached_proxy/cache"
	"cached_proxy/executor"
	"cached_proxy/feign"
	"log"
	"net/http"
	"time"
)

// 这里是初始化代码， 用于初始化各种服务

var (
	// Client 是爬虫服务的客户端
	Client feign.SpiderClient = feign.NewSpiderClientImpl(SpiderUrl, // SpiderUrl 是爬虫服务的地址， 通过环境变量 SPIDER_URL 设置
		http.Client{})
	// StudentService 是学生服务
	StudentService feign.StudentService = feign.NewStudentServiceImpl(&Client)
)

var (
	// AccountRepository 是账户的数据仓库
	AccountRepository = account.NewFileRepository("./_data")
	// AccountService 是账户的服务
	AccountService account.Service = account.NewServiceImpl(AccountRepository)
)

// updateTask 是一个通用的更新任务， 用于更新学生信息， 同时也会根据返回的错误信息进行账户锁定
func updateTask[V any](update func(*feign.Student) (*V, error)) func(string) (*V, bool) {
	return func(studentID string) (*V, bool) {
		student, err := StudentService.GetStudent(studentID)
		if err != nil {
			a, err := AccountService.GetAccountByAccountID(studentID)
			if err != nil {
				return nil, false
			}
			err = StudentService.SetStudent(a.AccountID(), a.GetPassword(), false)
			if err != nil {
				log.Printf("account %s is locked", studentID)
			}
		}
		if student == nil {
			return nil, false
		}
		value, err := update(student)
		if err != nil && err.Error() == "unauthorized" {
			log.Print("unauthorized: ", studentID)
			// 如果是未授权， 锁定账户
			err := AccountService.LockAccount(studentID)
			if err != nil {
				log.Print("failed to lock account: ", err)
			}
		}
		return value, err == nil
	}

}

var (
	PublicChecker    = cache.NewDailyStatusChecker[any](30 * time.Second)
	ClassroomChecker = PublicChecker
	CalendarChecker  = PublicChecker
)

var (
	PersonalChecker = cache.NewIntervalStatusChecker[any](2*time.Hour, 30*time.Second)
	InfoChecker     = PersonalChecker
	ScoreChecker    = PersonalChecker
	RankChecker     = PersonalChecker
	ExamChecker     = PersonalChecker
	CourseChecker   = PersonalChecker
)

var (
	TodayClassroomUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetClassroomStatus(0)
		return &value, err
	})
	TomorrowClassroomUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetClassroomStatus(1)
		return &value, err
	})
	CalendarUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetTeachingCalendar()
		return &value, err
	})
)

var (
	StudentInfoUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetInfo()
		return &value, err
	})
	StudentMajorScoreUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetStudentScore(true)
		return &value, err
	})
	StudentMinorScoreUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetStudentScore(false)
		return &value, err
	})
	StudentTotalRankUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetStudentRank(false)
		return &value, err
	})
	StudentRequiredRankUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetStudentRank(true)
		return &value, err
	})
	StudentExamUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetStudentExams()
		return &value, err
	})
	StudentCourseUpdater = updateTask[any](func(student *feign.Student) (*any, error) {
		value, err := (*student).GetStudentCourses()
		return &value, err
	})
)

var exec = executor.NewWorkerPool(10)

var (
	TodayClassroomService    = cache.NewPublicInformationService[any](exec, ClassroomChecker, TodayClassroomUpdater)
	TomorrowClassroomService = cache.NewPublicInformationService[any](exec, ClassroomChecker, TomorrowClassroomUpdater)
	CalendarService          = cache.NewPublicInformationService[any](exec, CalendarChecker, CalendarUpdater)
)

var (
	StudentInfoService         = cache.NewPersonalInformationService[any](exec, InfoChecker, StudentInfoUpdater)
	StudentMajorScoreService   = cache.NewPersonalInformationService[any](exec, ScoreChecker, StudentMajorScoreUpdater)
	StudentMinorScoreService   = cache.NewPersonalInformationService[any](exec, ScoreChecker, StudentMinorScoreUpdater)
	StudentTotalRankService    = cache.NewPersonalInformationService[any](exec, RankChecker, StudentTotalRankUpdater)
	StudentRequiredRankService = cache.NewPersonalInformationService[any](exec, RankChecker, StudentRequiredRankUpdater)
	StudentExamService         = cache.NewPersonalInformationService[any](exec, ExamChecker, StudentExamUpdater)
	StudentCourseService       = cache.NewPersonalInformationService[any](exec, CourseChecker, StudentCourseUpdater)
)
