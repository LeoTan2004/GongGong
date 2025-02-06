package main

import (
	"fmt"
	"net/http"
)

func StartApiServer(port int) {
	server := http.NewServeMux()
	server.HandleFunc("/login", Login)
	server.HandleFunc("/courses", CourseHandler.GetInfo)
	server.HandleFunc("/exams", ExamHandler.GetInfo)
	server.HandleFunc("/info", InfoHandler.GetInfo)
	server.HandleFunc("/scores", MajorScoreHandler.GetInfo)
	server.HandleFunc("/minor/scores", MinorScoreHandler.GetInfo)
	server.HandleFunc("/rank", TotalRankHandler.GetInfo)
	server.HandleFunc("/compulsory/rank", RequiredRankHandler.GetInfo)
	server.HandleFunc("/calendar", CalendarHandler.GetInfo)
	server.HandleFunc("/classroom/today", TodayClassroomHandler.GetInfo)
	server.HandleFunc("/classroom/tomorrow", TomorrowClassroomHandler.GetInfo)
	server.HandleFunc("/account", AccountHandler.GetInfo)
	server.HandleFunc("/icalendar/courses", CoursesCalendarHandler.GetInfo)
	server.HandleFunc("/icalendar/exams", ExamCalendarHandler.GetInfo)
	server.HandleFunc("/icalendar", CalPage)
	fmt.Printf("Proxy Server URL: %s\n", SpiderUrl)
	fmt.Printf("Starting server on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		return
	}
}

func main() {
	StartApiServer(ApiPort)
}
