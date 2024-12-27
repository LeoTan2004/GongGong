package cache

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试 HttpUpdater 中的 Send 方法
func TestHttpUpdater_Send(t *testing.T) {
	// 模拟的响应数据
	expectedResponse := ResponseData{
		Code: 0,
		Msg:  "success",
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	// 创建一个模拟的 HTTP 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证方法和头部信息
		if r.Method != "GET" {
			t.Errorf("期望的请求方法是 GET，实际是 %s", r.Method)
		}
		if r.Header.Get("token") != "test-token" {
			t.Errorf("期望的 token 是 'test-token'，实际是 '%s'", r.Header.Get("token"))
		}

		// 返回 JSON 响应
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	// 创建 HttpUpdater 实例
	updater := HttpUpdater{
		url:    server.URL,
		client: http.Client{},
	}

	// 调用 Send 方法
	headers := map[string]string{"token": "test-token"}
	data, err := updater.Send("GET", headers, nil)
	if err != nil {
		t.Fatalf("Send 方法返回错误: %v", err)
	}

	// 验证响应结果
	responseData, ok := data.(ResponseData)
	if !ok {
		t.Fatalf("响应结果的类型不正确: %T", data)
	}
	if responseData.Code != expectedResponse.Code {
		t.Errorf("期望的响应 Code 是 %d，实际是 %d", expectedResponse.Code, responseData.Code)
	}
	if responseData.Msg != expectedResponse.Msg {
		t.Errorf("期望的响应 Msg 是 '%s'，实际是 '%s'", expectedResponse.Msg, responseData.Msg)
	}
}
