package cache

import (
	"bytes"
	"cached_proxy/repo"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Updater 更新函数
type Updater[K, V any] interface {
	Invoke(key K) (V, error) // 更新执行入口
}

type HttpUpdater struct {
	url    string
	client http.Client
}

func (u *HttpUpdater) Send(method string, headers map[string]string, data interface{}) (res interface{}, err error) {
	var body io.Reader = nil
	if method == "POST" && data != nil {
		// 将数据编码为 JSON 格式
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("JSON 编码失败")
			return nil, err
		}
		bytes.NewBuffer(jsonData)
	}
	r, err := http.NewRequest(method, u.url, body)
	if err != nil {
		log.Println("创建请求失败")
		return nil, err
	}
	for key, value := range headers {
		r.Header.Set(key, value)
	}
	resp, err := u.client.Do(r)
	if err != nil {
		log.Printf("请求失败 【%s】\n", u.url)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("请求关闭失败-【%s】\n", u.url)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("请求失败，状态码: %d\n", resp.StatusCode)
		return "", nil
	}
	// 解码 JSON 响应
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("JSON 解码失败:", err)
		return nil, err
	}
	return result, nil
}

type HttpGetWithTokenUpdater[K string, V string] struct {
	HttpUpdater
	tokenRepo repo.KVRepo[string, string] // token存储器
}

func (h *HttpGetWithTokenUpdater[K, V]) Invoke(key string) (result V, err error) {
	token, found := h.tokenRepo.Get(key)
	if !found {
		log.Printf("账户不存在-【%s】\n", key)
		err = fmt.Errorf("账户不存在-【%s】\n", key)
		return "", err
	}
	// 发送请求
	headers := map[string]string{"token": token}
	data, err := h.Send("GET", headers, nil)
	if err != nil {
		return "", err
	}
	result, ok := data.(V)
	if !ok {
		return result, fmt.Errorf("类型断言失败，无法将 %T 转换为目标类型", data)
	}
	return result, nil
}
