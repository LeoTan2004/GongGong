package feign

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

// SpiderClientImpl 是 SpiderClient 接口的具体实现。
type SpiderClientImpl struct {
	baseUrl string
	client  http.Client
}

// buildRequest 构建实际请求
func (c *SpiderClientImpl) buildRequest(method string, uri string, token string, data any) (*http.Request, error) {
	// 参数合法性验证
	if method == "" || uri == "" {
		return nil, fmt.Errorf("method 和 uri 都不能为空")
	}

	headers := map[string]string{}

	if token != "" {
		headers["token"] = token
	}

	// 构造请求 URL
	u, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("解析 baseUrl 失败: %w", err)
	}
	u.Path = path.Join(u.Path, uri) // 保留 baseURL 的 host 和 scheme
	actualRequestUrl := u.String()

	// 初始化请求 body
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("JSON 编码失败 (method: %s, uri: %s): %v", method, uri, err)
			return nil, fmt.Errorf("JSON 编码失败: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
		headers["Content-Type"] = "application/json"
	}

	// 创建 HTTP 请求
	r, err := http.NewRequest(method, actualRequestUrl, body)
	if err != nil {
		log.Printf("创建请求失败 (method: %s, url: %s): %v", method, actualRequestUrl, err)
		return nil, fmt.Errorf("创建请求失败 (method: %s, url: %s): %w", method, actualRequestUrl, err)
	}

	// 设置请求头
	for key, value := range headers {
		r.Header.Set(key, value)
	}
	// 返回请求
	return r, nil
}

// 发送请求
func (c *SpiderClientImpl) sendRequest(r *http.Request) (*http.Response, error) {
	// 发起 HTTP 请求
	response, err := c.client.Do(r)
	if err != nil {
		log.Printf("请求失败: method=%s, url=%s, error=%v", r.Method, r.URL, err)
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}

	// 检查响应状态码
	if response.StatusCode != http.StatusOK {
		log.Printf("异常返回: method=%s, url=%s, status=%d", r.Method, r.URL, response.StatusCode)
	}

	return response, nil
}

// decodeResponse 解码统一返回
func (c *SpiderClientImpl) decodeResponse(response *http.Response) (CommonResponse[any], error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("返回解析失败: %v", err)
		}
	}(response.Body)
	var result CommonResponse[any]
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		fmt.Println("返回体 JSON 解码失败：", err)
		return CommonResponse[any]{}, err
	}
	return result, nil
}

// getWithToken 发送头部携带token的get请求
func (c *SpiderClientImpl) getWithToken(uri string, token string) (CommonResponse[any], error) {
	request, err := c.buildRequest("GET", uri, token, nil)
	if err != nil {
		return CommonResponse[any]{}, err
	}
	response, err := c.sendRequest(request)
	if err != nil {
		return CommonResponse[any]{}, err
	}
	commonResponse, err := c.decodeResponse(response)
	if err != nil {
		return CommonResponse[any]{}, err
	}
	if commonResponse.Code != 1 {
		return CommonResponse[any]{}, fmt.Errorf("返回错误：%d，返回信息：%s", commonResponse.Code, commonResponse.Message)
	}
	return commonResponse, nil
}

func (c *SpiderClientImpl) GetTeachingCalendar(token string) (any, error) {
	commonResponse, err := c.getWithToken("/calendar", token)
	if err != nil {
		return nil, err
	}
	return commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetClassroomStatus(token string, day int) (any, error) {

	var commonResponse CommonResponse[any]
	var err error
	if day == 0 {
		commonResponse, err = c.getWithToken("/classroom/today", token)
	} else if day == 1 {
		commonResponse, err = c.getWithToken("/classroom/tomorrow", token)
	} else {
		return nil, fmt.Errorf("day只能为0或1")
	}
	if err != nil {
		return nil, err
	}
	return commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentCourses(token string) (any, error) {
	commonResponse, err := c.getWithToken("/courses", token)
	if err != nil {
		return nil, err
	}
	return commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentExams(token string) (any, error) {
	commonResponse, err := c.getWithToken("/exams", token)
	if err != nil {
		return nil, err
	}
	return commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentInfo(token string) (any, error) {
	commonResponse, err := c.getWithToken("/info", token)
	if err != nil {
		return nil, err
	}
	return commonResponse.Data, nil
}

func (c *SpiderClientImpl) Login(username string, password string) (LoginResponse, error) {
	request, err := c.buildRequest("POST", "/login", "", map[string]string{"username": username, "password": password})
	if err != nil {
		return LoginResponse{}, err
	}
	response, err := c.sendRequest(request)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return LoginResponse{}, err
	}
	commonResponse, err := c.decodeResponse(response)
	if err != nil {
		return LoginResponse{}, err
	}
	// 解析数据
	dataMap, ok := commonResponse.Data.(map[string]interface{})
	if !ok {
		return LoginResponse{}, fmt.Errorf("解析返回数据失败: %v", commonResponse.Data)
	}

	// 转换为 LoginResponse
	token, ok := dataMap["token"].(string)
	if !ok {
		return LoginResponse{}, fmt.Errorf("返回数据中缺少 token")
	}

	return LoginResponse{Token: token}, nil
}

func (c *SpiderClientImpl) GetStudentScore(token string, isMajor bool) (any, error) {
	var commonResponse CommonResponse[any]
	var err error
	if isMajor {
		commonResponse, err = c.getWithToken("/scores", token)
	} else {
		commonResponse, err = c.getWithToken("/minor/scores", token)
	}
	if err != nil {
		return nil, err
	}
	return commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentRank(token string, onlyRequired bool) (any, error) {
	var commonResponse CommonResponse[any]
	var err error
	if onlyRequired {
		// TODO(2024年12月28日 11:04 , LeoTan) 添加仅仅必修课程的排名计算的接口
		return nil, fmt.Errorf("onlyRequired is not supported")
		//commonResponse, err = c.getWithToken("/rank", token)
	} else {
		commonResponse, err = c.getWithToken("/rank", token)
	}
	if err != nil {
		return nil, err
	}
	return commonResponse.Data, nil
}

func NewSpiderClientImpl(baseUrl string, client http.Client) *SpiderClientImpl {
	if client.Timeout == 0 {
		client.Timeout = 5 * time.Second
	}
	return &SpiderClientImpl{baseUrl: baseUrl, client: client}
}
