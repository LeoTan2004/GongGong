package main

import (
	account2 "cached_proxy/account"
	"cached_proxy/cache"
	"cached_proxy/feign"
	"encoding/json"
	"log"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = StudentService.SetStudent(creds.Username, creds.Password, true)
	if err == nil {
		token, err := AccountService.Login(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		resp := feign.CommonResponse[map[string]string]{
			Code:    1,
			Message: "success",
			Data: map[string]string{
				"token": token,
			},
		}
		// 返回 token
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}
	if err.Error() == "unauthorized" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type InfoGetter[V any] struct {
	info cache.InformationService[V]
}

func (c *InfoGetter[V]) GetInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetInfo %s\n", r.RequestURI)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.Header.Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	account, err := AccountService.GetAccountByToken(token)
	if err != nil {
		if err.Error() == "account not found" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	if account == nil || account.Status() != account2.Normal {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	info, err := c.info.GetInfo(account.AccountID())
	if err != nil {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
	}
	resp := feign.CommonResponse[any]{
		Code:    1,
		Message: "success",
		Data:    info,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

var (
	CalendarHandler          = &InfoGetter[any]{info: CalendarService}
	TodayClassroomHandler    = &InfoGetter[any]{info: TodayClassroomService}
	TomorrowClassroomHandler = &InfoGetter[any]{info: TomorrowClassroomService}
	InfoHandler              = &InfoGetter[any]{info: StudentInfoService}
	MajorScoreHandler        = &InfoGetter[any]{info: StudentMajorScoreService}
	MinorScoreHandler        = &InfoGetter[any]{info: StudentMinorScoreService}
	TotalRankHandler         = &InfoGetter[any]{info: StudentTotalRankService}
	RequiredRankHandler      = &InfoGetter[any]{info: StudentRequiredRankService}
	ExamHandler              = &InfoGetter[any]{info: StudentExamService}
	CourseHandler            = &InfoGetter[any]{info: StudentCourseService}
)
