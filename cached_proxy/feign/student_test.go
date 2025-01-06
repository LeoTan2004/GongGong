package feign

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

type MockSpiderClient struct {
	expectedLoginResponse LoginResponse
	expectedReturn        any
	expectedError         error
}

func (m *MockSpiderClient) GetTeachingCalendar(_ string) (any, error) {
	return m.expectedReturn, m.expectedError
}

func (m *MockSpiderClient) GetClassroomStatus(_ string, _ int) (any, error) {
	return m.expectedReturn, m.expectedError
}

func (m *MockSpiderClient) GetStudentCourses(_ string) (any, error) {
	return m.expectedReturn, m.expectedError
}

func (m *MockSpiderClient) GetStudentExams(_ string) (any, error) {
	return m.expectedReturn, m.expectedError
}

func (m *MockSpiderClient) GetStudentInfo(_ string) (any, error) {
	return m.expectedReturn, m.expectedError
}

func (m *MockSpiderClient) Login(_ string, _ string) (LoginResponse, error) {
	return m.expectedLoginResponse, m.expectedError
}

func (m *MockSpiderClient) GetStudentScore(_ string, _ bool) (any, error) {
	return m.expectedReturn, m.expectedError
}

func (m *MockSpiderClient) GetStudentRank(_ string, _ bool) (any, error) {
	return m.expectedReturn, m.expectedError
}

// NewStudent 创建一个新的学生账户。
func (m *MockSpiderClient) NewStudent(username string, password string) (Student, error) {
	return &StudentImpl{
		username: username,
		password: password,
	}, nil
}

func TestStudentImpl_doGetter(t *testing.T) {
	retryTime := 3
	tests := []struct {
		name              string
		token             string
		mockFunction      func(token string) (any, error)
		mockError         error
		mockLoginResponse LoginResponse
		expectedError     bool
		expectedErrorMsg  string
		expectedResult    any
	}{
		{
			name:  "doGetter success",
			token: "valid-token",
			mockFunction: func(token string) (any, error) {
				return "valid data", nil
			},
			expectedError:  false,
			expectedResult: "valid data",
		},
		{
			name:  "doGetter success invalid token and retry login success",
			token: "invalid-token",
			mockFunction: func(token string) (any, error) {
				if token != "valid-token" {
					return nil, errors.New("unauthorized")
				} else {
					return "valid data", nil
				}
			},
			mockLoginResponse: LoginResponse{Token: "valid-token"},
			expectedError:     false,
			expectedResult:    "valid data",
		},
		{
			name:  "doGetter success valid token and retry success",
			token: "invalid-token",
			mockFunction: func(token string) (any, error) {
				if retryTime == 1 {
					return "valid data", nil
				} else {
					retryTime--
					return nil, errors.New("service unavailable")
				}
			},
			mockError:      nil,
			expectedError:  false,
			expectedResult: "valid data",
		}, {
			name: "doGetter failure unauthorized username and password",
			mockFunction: func(token string) (any, error) {
				return nil, errors.New("unauthorized")
			},
			mockError:        fmt.Errorf("unauthorized"),
			expectedError:    true,
			expectedErrorMsg: "unauthorized",
			expectedResult:   nil,
		},
		{
			name:  "doGetter failure exceeded retry attempts with service unavailable",
			token: "any-token",
			mockFunction: func(token string) (any, error) {
				return nil, errors.New("service unavailable")
			},
			expectedError:    true,
			expectedErrorMsg: "exceeded retry attempts: service unavailable",
			expectedResult:   nil,
		}, {
			name: "doGetter failure with unknown error",
			mockFunction: func(token string) (any, error) {
				return nil, errors.New("unknown error")
			},
			expectedError:    true,
			expectedErrorMsg: "unknown error",
			expectedResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			student := &StudentImpl{
				dynamicToken: tt.token,
				spider: &MockSpiderClient{
					expectedReturn:        tt.expectedResult,
					expectedError:         tt.mockError,
					expectedLoginResponse: tt.mockLoginResponse,
				},
			}

			result, err := student.doGetter(tt.mockFunction)
			if err != nil && err.Error() != tt.expectedErrorMsg {
				t.Fatalf("Expected error message: %v, got: %v", tt.expectedErrorMsg, err.Error())
			}
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if !tt.expectedError && result != tt.expectedResult {
				t.Errorf("Expected result: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// UUID format: xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
	uuidRegex := `^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`
	matched, _ := regexp.MatchString(uuidRegex, uuid)
	if !matched {
		t.Fatalf("Generated UUID %s does not match the expected format", uuid)
	}
}

func TestGenerateMultipleUUIDs(t *testing.T) {
	uuidSet := make(map[string]struct{})
	count := 1000

	for i := 0; i < count; i++ {
		uuid, err := GenerateUUID()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		// Check for uniqueness
		if _, exists := uuidSet[uuid]; exists {
			t.Fatalf("Duplicate UUID found: %s", uuid)
		}
		uuidSet[uuid] = struct{}{}
	}
}

func TestNewStudentImpl(t *testing.T) {
	type args struct {
		username string
		password string
		client   SpiderClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "NewStudentImpl success",
			args: args{
				username: "valid-username",
				password: "valid-password",
				client: &MockSpiderClient{
					expectedLoginResponse: LoginResponse{Token: "valid-token"},
					expectedError:         nil,
					expectedReturn:        "valid data",
				},
			},
			wantErr: false,
		},
		{
			name: "NewStudentImpl failure with invalid username and password",
			args: args{
				username: "invalid-username",
				password: "invalid-password",
				client: &MockSpiderClient{
					expectedError: fmt.Errorf("unauthorized"),
				},
			},
			wantErr: true,
		},
		{
			name: "NewStudentImpl failure with service unavailable",
			args: args{
				username: "valid-username",
				password: "valid-password",
				client: &MockSpiderClient{
					expectedError: fmt.Errorf("service unavailable"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStudentImpl(tt.args.username, tt.args.password, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStudentImpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStudentImplIntegrate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// login request
		if r.Method == "POST" && r.URL.Path == "/login" {
			var actualBody map[string]string
			if err := json.NewDecoder(r.Body).Decode(&actualBody); err != nil {
				t.Fatalf("Failed to decode request body: %v", err)
			}
			if actualBody["username"] != "valid-user" || actualBody["password"] != "valid-password" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"code": 1,
				"message": "success",
				"data": {
					"token": "valid-token"
				}
			}`))
			return
		}
		// get request
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		token := r.Header.Get("token")
		if token != "valid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"code": 1,
			"message": "success",
			"data": {
				"data": "valid-data"
			}
		}`))
	}))
	defer server.Close()

	validUsername := "valid-user"
	validPassword := "valid-password"
	baseUrl := server.URL
	client := NewSpiderClientImpl(baseUrl, http.Client{})
	var student Student
	var err error
	var data any
	// Test NewStudentImpl
	t.Run("NewStudentImpl success with valid username and password", func(t *testing.T) {
		student, err = NewStudentImpl(validUsername, validPassword, client)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("NewStudentImpl failure with invalid username and password", func(t *testing.T) {
		_, err = NewStudentImpl("invalid-username", "invalid-password", client)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Test GetTeachingCalendar
	t.Run("GetTeachingCalendar success", func(t *testing.T) {
		data, err = student.GetTeachingCalendar()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetTeachingCalendar success with invalid token and retry login", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     validUsername,
			password:     validPassword,
			dynamicToken: "invalid-token",
			spider:       client,
		}
		data, err = invalidTokenStudent.GetTeachingCalendar()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetTeachingCalendar failure with invalid token and retry login failure", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     "invalid-username",
			password:     "invalid-password",
			dynamicToken: "invalid",
			spider:       client,
		}
		_, err = invalidTokenStudent.GetTeachingCalendar()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Test GetClassroomStatus
	t.Run("GetClassroomStatus success", func(t *testing.T) {
		data, err = student.GetClassroomStatus(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetClassroomStatus success with invalid token and retry login", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     validUsername,
			password:     validPassword,
			dynamicToken: "invalid",
			spider:       client,
		}
		data, err = invalidTokenStudent.GetClassroomStatus(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetClassroomStatus failure with invalid token and retry login failure", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     "invalid-username",
			password:     "invalid-password",
			dynamicToken: "invalid-token",
			spider:       client,
		}
		_, err = invalidTokenStudent.GetClassroomStatus(1)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// GetStudentCourses
	t.Run("GetStudentCourses success", func(t *testing.T) {
		data, err = student.GetStudentCourses()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentCourses success with invalid token and retry login", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     validUsername,
			password:     validPassword,
			dynamicToken: "invalid",
			spider:       client,
		}
		data, err = invalidTokenStudent.GetStudentCourses()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentCourses failure with invalid token and retry login failure", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     "invalid-username",
			password:     "invalid-password",
			dynamicToken: "invalid",
			spider:       client,
		}
		_, err = invalidTokenStudent.GetStudentCourses()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// GetStudentExams
	t.Run("GetStudentExams success", func(t *testing.T) {
		data, err = student.GetStudentExams()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentExams success with invalid token and retry login", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     validUsername,
			password:     validPassword,
			dynamicToken: "invalid",
			spider:       client,
		}
		data, err = invalidTokenStudent.GetStudentExams()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentExams failure with invalid token and retry login failure", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     "invalid-username",
			password:     "invalid-password",
			dynamicToken: "invalid",
			spider:       client,
		}
		_, err = invalidTokenStudent.GetStudentExams()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// GetStudentInfo
	t.Run("GetStudentInfo success", func(t *testing.T) {
		data, err = student.GetInfo()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentInfo success with invalid token and retry login", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     validUsername,
			password:     validPassword,
			dynamicToken: "invalid",
			spider:       client,
		}
		data, err = invalidTokenStudent.GetInfo()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentInfo failure with invalid token and retry login failure", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     "invalid-username",
			password:     "invalid-password",
			dynamicToken: "invalid",
			spider:       client,
		}
		_, err = invalidTokenStudent.GetInfo()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// GetStudentScore
	t.Run("GetStudentScore success", func(t *testing.T) {
		data, err = student.GetStudentScore(true)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
		data, err = student.GetStudentScore(false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentScore success with invalid token and retry login", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     validUsername,
			password:     validPassword,
			dynamicToken: "invalid",
			spider:       client,
		}
		data, err = invalidTokenStudent.GetStudentScore(true)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
		data, err = invalidTokenStudent.GetStudentScore(false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentScore failure with invalid token and retry login failure", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     "invalid-username",
			password:     "invalid-password",
			dynamicToken: "invalid",
			spider:       client,
		}
		_, err = invalidTokenStudent.GetStudentScore(true)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		_, err = invalidTokenStudent.GetStudentScore(false)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// GetStudentRank
	t.Run("GetStudentRank success", func(t *testing.T) {
		data, err = student.GetStudentRank(true)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
		data, err = student.GetStudentRank(false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentRank success with invalid token and retry login", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     validUsername,
			password:     validPassword,
			dynamicToken: "invalid",
			spider:       client,
		}
		data, err = invalidTokenStudent.GetStudentRank(true)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
		data, err = invalidTokenStudent.GetStudentRank(false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if data == "" {
			t.Fatalf("Expected data, got nil")
		}
	})

	t.Run("GetStudentRank failure with invalid token and retry login failure", func(t *testing.T) {
		var invalidTokenStudent Student = &StudentImpl{
			username:     "invalid-username",
			password:     "invalid-password",
			dynamicToken: "invalid",
			spider:       client,
		}
		_, err = invalidTokenStudent.GetStudentRank(true)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		_, err = invalidTokenStudent.GetStudentRank(false)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}
