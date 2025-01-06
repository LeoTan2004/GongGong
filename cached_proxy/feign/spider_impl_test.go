package feign

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpiderClient_Login(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		password       string
		mockResponse   string
		mockStatusCode int
		expectedError  bool
		expectedToken  string
	}{
		{
			name:     "Login success",
			username: "valid-user",
			password: "valid-password",
			mockResponse: `{
				"code": 1,
				"message": "success",
				"data": {
					"token": "a685e58e-1040-43e7-8c9c-5b2e3c0e7ec3"
				}
			}`,
			mockStatusCode: http.StatusOK,
			expectedError:  false,
			expectedToken:  "a685e58e-1040-43e7-8c9c-5b2e3c0e7ec3",
		},
		{
			name:           "Login failure - incorrect credentials",
			username:       "invalid-user",
			password:       "wrong-password",
			mockResponse:   ``,
			mockStatusCode: http.StatusUnauthorized,
			expectedError:  true,
			expectedToken:  "",
		},
		{
			name:           "Login failure - account not initialized",
			username:       "uninitialized-user",
			password:       "any-password",
			mockResponse:   ``,
			mockStatusCode: http.StatusConflict,
			expectedError:  true,
			expectedToken:  "",
		},
		{
			name:           "Login failure - system timeout",
			username:       "any-user",
			password:       "any-password",
			mockResponse:   ``,
			mockStatusCode: http.StatusServiceUnavailable,
			expectedError:  true,
			expectedToken:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected method POST, got %s", r.Method)
				}
				if r.URL.Path != "/login" {
					t.Errorf("Expected URL path /login, got %s", r.URL.Path)
				}

				var actualBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&actualBody); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}
				if actualBody["username"] != tt.username || actualBody["password"] != tt.password {
					t.Errorf("Expected username: %s, password: %s, got username: %s, password: %s",
						tt.username, tt.password, actualBody["username"], actualBody["password"])
				}

				w.WriteHeader(tt.mockStatusCode)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := NewSpiderClientImpl(server.URL, http.Client{})

			response, err := client.Login(tt.username, tt.password)

			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if !tt.expectedError && response.Token != tt.expectedToken {
				t.Errorf("Expected token: %s, got: %s", tt.expectedToken, response.Token)
			}
		})
	}
}

func TestSpiderClient_getWithToken(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		mockResponse   string
		mockStatusCode int
		expectedError  bool
		expectedData   any
	}{
		{
			name:  "GetWithToken success",
			token: "valid-token",
			mockResponse: `{
					"code": 1,
					"message": "success",
					"data": "truly data"

				}`,
			mockStatusCode: http.StatusOK,
			expectedError:  false,
			expectedData:   "truly data",
		},
		{
			name:           "GetWithToken failure - invalid token",
			token:          "invalid-token",
			mockResponse:   ``,
			mockStatusCode: http.StatusUnauthorized,
			expectedError:  true,
			expectedData:   nil,
		},
		{
			name:           "GetWithToken failure - system timeout",
			token:          "any-token",
			mockResponse:   ``,
			mockStatusCode: http.StatusServiceUnavailable,
			expectedError:  true,
			expectedData:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected method GET, got %s", r.Method)
				}
				if r.URL.Path != "/test-uri" {
					t.Errorf("Expected URL path /calendar, got %s", r.URL.Path)
				}
				if r.Header.Get("token") != tt.token {
					t.Errorf("Expected token: %s, got: %s", tt.token, r.Header.Get("token"))
				}

				w.WriteHeader(tt.mockStatusCode)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := NewSpiderClientImpl(server.URL, http.Client{})

			response, err := client.getWithToken("/test-uri", tt.token)

			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if !tt.expectedError && response.Data != tt.expectedData {
				t.Errorf("Expected data: %v, got: %v", tt.expectedData, response.Data)
			}
		})
	}
}

func TestSpiderClientImplIntegrate(t *testing.T) {
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
	validToken := "valid-token"
	baseUrl := server.URL
	client := NewSpiderClientImpl(baseUrl, http.Client{})
	t.Run("Integrate Test Login With Valid Username and password", func(t *testing.T) {
		login, err := client.Login(validUsername, validPassword)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		validToken = login.Token
	})
	login, err := client.Login(validUsername, validPassword)
	if err != nil {
		t.Fatalf("Expected error: %v, got: %v", nil, err)
	}
	validToken = login.Token

	t.Run("Integrate Test Login With Invalid Username and password", func(t *testing.T) {
		_, err := client.Login("invalid-user", "invalid-password")
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
	})

	t.Run("Integrate Test GetClassroomStatus With Valid Token", func(t *testing.T) {
		resp, err := client.GetClassroomStatus(validToken, 0)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
		resp, err = client.GetClassroomStatus(validToken, 1)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
	})

	t.Run("Integrate Test GetClassroomStatus With Invalid Token", func(t *testing.T) {
		_, err := client.GetClassroomStatus("invalid-token", 0)
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
		_, err = client.GetClassroomStatus("invalid-token", 1)
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
	})

	t.Run("Integrate Test GetStudentInfo With Valid Token", func(t *testing.T) {
		resp, err := client.GetStudentInfo(validToken)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
	})

	t.Run("Integrate Test GetStudentInfo With Invalid Token", func(t *testing.T) {
		_, err := client.GetStudentInfo("invalid-token")
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
	})

	t.Run("Integrate Test GetStudentCourses With Valid Token", func(t *testing.T) {
		resp, err := client.GetStudentCourses(validToken)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
	})

	t.Run("Integrate Test GetStudentCourses With Invalid Token", func(t *testing.T) {
		_, err := client.GetStudentCourses("invalid-token")
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
	})

	t.Run("Integrate Test GetStudentExams With Valid Token", func(t *testing.T) {
		resp, err := client.GetStudentExams(validToken)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
	})

	t.Run("Integrate Test GetStudentExams With Invalid Token", func(t *testing.T) {
		_, err := client.GetStudentExams("invalid-token")
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
	})

	t.Run("Integrate Test GetTeachingCalendar With Valid Token", func(t *testing.T) {
		resp, err := client.GetTeachingCalendar(validToken)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
	})

	t.Run("Integrate Test GetTeachingCalendar With Invalid Token", func(t *testing.T) {
		_, err := client.GetTeachingCalendar("invalid-token")
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
	})

	t.Run("Integrate Test GetStudentRank With Valid Token", func(t *testing.T) {
		resp, err := client.GetStudentRank(validToken, false)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
		resp, err = client.GetStudentRank(validToken, true)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}

	})

	t.Run("Integrate Test GetStudentRank With Invalid Token", func(t *testing.T) {
		_, err := client.GetStudentRank("invalid-token", false)
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
		_, err = client.GetStudentRank("invalid-token", true)
		if err == nil || err.Error() != "unauthorized" {
			t.Fatalf("Expected error: %v, got: %v", "unauthorized", err)
		}
	})

	t.Run("Integrate Test GetStudentScore With Valid Token", func(t *testing.T) {
		resp, err := client.GetStudentScore(validToken, false)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}
		resp, err = client.GetStudentScore(validToken, true)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v", nil, err)
		}
		if resp == "" {
			t.Fatalf("Expected data is not nil, but got: %v", resp)
		}

	})

}
