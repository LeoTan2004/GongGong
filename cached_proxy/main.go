package main

import (
	"fmt"
	"net/http"
)

func main() {

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
	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", server)
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		return
	}

}
