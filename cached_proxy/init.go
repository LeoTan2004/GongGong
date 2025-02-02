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
	ClassroomChecker = cache.NewDailyStatusChecker[feign.ClassroomStatusTable](30 * time.Second)
	CalendarChecker  = cache.NewDailyStatusChecker[feign.TeachingCalendar](30 * time.Second)
)

var (
	InfoChecker   = cache.NewIntervalStatusChecker[feign.StudentInfo](2*time.Hour, 30*time.Second)
	ScoreChecker  = cache.NewIntervalStatusChecker[feign.ScoreBoard](2*time.Hour, 30*time.Second)
	RankChecker   = cache.NewIntervalStatusChecker[feign.Rank](2*time.Hour, 30*time.Second)
	ExamChecker   = cache.NewIntervalStatusChecker[feign.ExamList](2*time.Hour, 30*time.Second)
	CourseChecker = cache.NewIntervalStatusChecker[feign.CourseList](2*time.Hour, 30*time.Second)
)

var (
	TodayClassroomUpdater = updateTask[feign.ClassroomStatusTable](func(student *feign.Student) (*feign.ClassroomStatusTable, error) {
		value, err := (*student).GetClassroomStatus(0)
		return value, err
	})
	TomorrowClassroomUpdater = updateTask[feign.ClassroomStatusTable](func(student *feign.Student) (*feign.ClassroomStatusTable, error) {
		value, err := (*student).GetClassroomStatus(1)
		return value, err
	})
	CalendarUpdater = updateTask[feign.TeachingCalendar](func(student *feign.Student) (*feign.TeachingCalendar, error) {
		value, err := (*student).GetTeachingCalendar()
		return value, err
	})
)

var (
	StudentInfoUpdater = updateTask[feign.StudentInfo](func(student *feign.Student) (*feign.StudentInfo, error) {
		value, err := (*student).GetInfo()
		return value, err
	})
	StudentMajorScoreUpdater = updateTask[feign.ScoreBoard](func(student *feign.Student) (*feign.ScoreBoard, error) {
		value, err := (*student).GetStudentScore(true)
		return value, err
	})
	StudentMinorScoreUpdater = updateTask[feign.ScoreBoard](func(student *feign.Student) (*feign.ScoreBoard, error) {
		value, err := (*student).GetStudentScore(false)
		return value, err
	})
	StudentTotalRankUpdater = updateTask[feign.Rank](func(student *feign.Student) (*feign.Rank, error) {
		value, err := (*student).GetStudentRank(false)
		return value, err
	})
	StudentRequiredRankUpdater = updateTask[feign.Rank](func(student *feign.Student) (*feign.Rank, error) {
		value, err := (*student).GetStudentRank(true)
		return value, err
	})
	StudentExamUpdater = updateTask[feign.ExamList](func(student *feign.Student) (*feign.ExamList, error) {
		value, err := (*student).GetStudentExams()
		return value, err
	})
	StudentCourseUpdater = updateTask[feign.CourseList](func(student *feign.Student) (*feign.CourseList, error) {
		value, err := (*student).GetStudentCourses()
		return value, err
	})
)

var exec = executor.NewWorkerPool(10)

var (
	TodayClassroomService    = cache.NewPublicInformationService[feign.ClassroomStatusTable](exec, ClassroomChecker, TodayClassroomUpdater)
	TomorrowClassroomService = cache.NewPublicInformationService[feign.ClassroomStatusTable](exec, ClassroomChecker, TomorrowClassroomUpdater)
	CalendarService          = cache.NewPublicInformationService[feign.TeachingCalendar](exec, CalendarChecker, CalendarUpdater)
)

var (
	StudentInfoService         = cache.NewPersonalInformationService[feign.StudentInfo](exec, InfoChecker, StudentInfoUpdater)
	StudentMajorScoreService   = cache.NewPersonalInformationService[feign.ScoreBoard](exec, ScoreChecker, StudentMajorScoreUpdater)
	StudentMinorScoreService   = cache.NewPersonalInformationService[feign.ScoreBoard](exec, ScoreChecker, StudentMinorScoreUpdater)
	StudentTotalRankService    = cache.NewPersonalInformationService[feign.Rank](exec, RankChecker, StudentTotalRankUpdater)
	StudentRequiredRankService = cache.NewPersonalInformationService[feign.Rank](exec, RankChecker, StudentRequiredRankUpdater)
	StudentExamService         = cache.NewPersonalInformationService[feign.ExamList](exec, ExamChecker, StudentExamUpdater)
	StudentCourseService       = cache.NewPersonalInformationService[feign.CourseList](exec, CourseChecker, StudentCourseUpdater)
)
