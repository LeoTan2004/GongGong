package feign

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpiderClient_GetTeachingCalendar(t *testing.T) {
	runGetTests(t, "/calendar", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetTeachingCalendar(token)
	})
}

func TestSpiderClient_GetClassroomStatus(t *testing.T) {
	runGetTests(t, "/classroom/today", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetClassroomStatus(token, 0)
	})
	runGetTests(t, "/classroom/tomorrow", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetClassroomStatus(token, 1)
	})
}

func TestSpiderClient_GetStudentCourses(t *testing.T) {
	runGetTests(t, "/courses", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetStudentCourses(token)
	})
}

func TestSpiderClient_GetStudentExams(t *testing.T) {
	runGetTests(t, "/exams", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetStudentExams(token)
	})
}

func TestSpiderClient_GetStudentInfo(t *testing.T) {
	runGetTests(t, "/info", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetStudentInfo(token)
	})
}

func TestSpiderClient_GetStudentScore(t *testing.T) {
	runGetTests(t, "/scores", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetStudentScore(token, true)
	})
	runGetTests(t, "/minor/scores", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetStudentScore(token, false)
	})
}

func TestSpiderClient_GetStudentRank(t *testing.T) {
	runGetTests(t, "/rank", func(client *SpiderClientImpl, token string) (any, error) {
		return client.GetStudentRank(token, false)
	})
}

// runGetTests 通用的测试逻辑
func runGetTests(t *testing.T, uri string, methodToTest func(*SpiderClientImpl, string) (any, error)) {
	tests := []struct {
		name           string
		token          string
		mockResponse   string
		mockStatusCode int
		expectedError  bool
		expectedData   map[string]any
	}{
		{
			name:  "Valid token, successful response",
			token: "valid-token",
			mockResponse: `{
				"code": 1,
				"message": "success",
				"data": {
					"key": "value"
				}
			}`,
			mockStatusCode: http.StatusOK,
			expectedError:  false,
			expectedData: map[string]any{
				"key": "value",
			},
		},
		{
			name:  "Invalid token, error response",
			token: "invalid-token",
			mockResponse: `{
				"code": 0,
				"message": "invalid token"
			}`,
			mockStatusCode: http.StatusOK,
			expectedError:  true,
			expectedData:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟服务器响应
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 验证 token 是否正确
				if r.Header.Get("token") != tt.token {
					t.Errorf("Expected token %s, got %s", tt.token, r.Header.Get("token"))
				}
				// 验证路径是否正确
				if r.URL.Path != uri {
					t.Errorf("Expected URI %s, got %s", uri, r.URL.Path)
				}
				// 返回模拟响应
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// 创建 SpiderClient 实例
			client := NewSpiderClientImpl(server.URL, http.Client{})

			// 调用方法
			data, err := methodToTest(client, tt.token)

			// 验证是否符合预期
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if err == nil {
				// 比较数据内容
				if dataMap, ok := data.(map[string]any); ok {
					for key, expectedValue := range tt.expectedData {
						if dataMap[key] != expectedValue {
							t.Errorf("For key %s, expected %v, got %v", key, expectedValue, dataMap[key])
						}
					}
				} else {
					t.Errorf("Expected map[string]any, got %T", data)
				}
			}
		})
	}
}

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
			name:     "Login failure - incorrect credentials",
			username: "invalid-user",
			password: "wrong-password",
			mockResponse: `{
				"code": 0,
				"message": "账户密码错误",
				"data": null
			}`,
			mockStatusCode: http.StatusOK,
			expectedError:  true,
			expectedToken:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟服务器响应
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 验证请求方法和路径
				if r.Method != "POST" {
					t.Errorf("Expected method POST, got %s", r.Method)
				}
				if r.URL.Path != "/login" {
					t.Errorf("Expected URL path /login, got %s", r.URL.Path)
				}

				// 验证请求体
				var actualBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&actualBody); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}
				if actualBody["username"] != tt.username || actualBody["password"] != tt.password {
					t.Errorf("Expected username: %s, password: %s, got username: %s, password: %s",
						tt.username, tt.password, actualBody["username"], actualBody["password"])
				}

				// 返回模拟响应
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// 创建 SpiderClient 实例
			client := NewSpiderClientImpl(server.URL, http.Client{})

			// 调用 Login 方法
			response, err := client.Login(tt.username, tt.password)

			// 验证结果
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if !tt.expectedError && response.Token != tt.expectedToken {
				t.Errorf("Expected token: %s, got: %s", tt.expectedToken, response.Token)
			}
		})
	}
}
