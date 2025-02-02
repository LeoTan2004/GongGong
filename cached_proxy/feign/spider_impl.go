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
	switch response.StatusCode {
	case http.StatusOK:
		// 正常返回
		return response, nil
	case http.StatusUnauthorized:
		// 未授权
		log.Printf("Unauthorized: method=%s, url=%s", r.Method, r.URL)
		return nil, fmt.Errorf("unauthorized")
	case http.StatusServiceUnavailable:
		// 服务不可用
		log.Printf("ServiceUnavailable: method=%s, url=%s", r.Method, r.URL)
		return nil, fmt.Errorf("service unavailable")
	default:
		// 其他错误
		log.Printf("Unkown Error: method=%s, url=%s, status=%d", r.Method, r.URL, response.StatusCode)
		return nil, fmt.Errorf("unkown error: status=%d", response.StatusCode)
	}
}

// decodeResponse 解码统一返回
func decodeResponse[V any](response *http.Response) (CommonResponse[V], error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("返回解析失败: %v", err)
		}
	}(response.Body)
	var result CommonResponse[V]
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		fmt.Println("返回体 JSON 解码失败：", err)
		return CommonResponse[V]{}, err
	}
	return result, nil
}

// getWithToken 发送头部携带token的get请求
func getWithToken[V any](c *SpiderClientImpl, uri string, token string) (*CommonResponse[V], error) {
	// 构建请求
	request, err := c.buildRequest("GET", uri, token, nil)
	if err != nil {
		return nil, err
	}
	// 发送请求
	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}
	// 解析返回
	commonResponse, err := decodeResponse[V](response)
	if err != nil {
		return nil, err
	}
	if commonResponse.Code != 1 {
		return nil, fmt.Errorf("返回错误：%d，返回信息：%s", commonResponse.Code, commonResponse.Message)
	}
	return &commonResponse, nil
}

func (c *SpiderClientImpl) GetTeachingCalendar(token string) (*TeachingCalendar, error) {
	commonResponse, err := getWithToken[TeachingCalendar](c, "/calendar", token)
	if err != nil {
		return nil, err
	}
	return &commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetClassroomStatus(token string, day int) (*ClassroomStatusTable, error) {
	uri := fmt.Sprintf("/classroom/%d", day)
	commonResponse, err := getWithToken[ClassroomStatusTable](c, uri, token)
	if err != nil {
		return nil, err
	}
	return &commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentCourses(token string) (*CourseList, error) {
	commonResponse, err := getWithToken[CourseList](c, "/courses", token)
	if err != nil {
		return nil, err
	}
	return &commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentExams(token string) (*ExamList, error) {
	commonResponse, err := getWithToken[ExamList](c, "/exams", token)
	if err != nil {
		return nil, err
	}
	return &commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentInfo(token string) (*StudentInfo, error) {
	commonResponse, err := getWithToken[StudentInfo](c, "/info", token)
	if err != nil {
		return nil, err
	}
	return &commonResponse.Data, nil
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
	loginResponse, err := decodeResponse[LoginResponse](response)
	if err != nil {
		return LoginResponse{}, err
	}
	return loginResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentScore(token string, isMajor bool) (*ScoreBoard, error) {
	var commonResponse *CommonResponse[ScoreBoard]
	var err error
	if isMajor {
		commonResponse, err = getWithToken[ScoreBoard](c, "/scores", token)
	} else {
		commonResponse, err = getWithToken[ScoreBoard](c, "/minor/scores", token)
	}
	if err != nil {
		return nil, err
	}
	return &commonResponse.Data, nil
}

func (c *SpiderClientImpl) GetStudentRank(token string, onlyRequired bool) (*Rank, error) {
	var commonResponse *CommonResponse[Rank]
	var err error
	if onlyRequired {
		commonResponse, err = getWithToken[Rank](c, "/compulsory/rank", token)
	} else {
		commonResponse, err = getWithToken[Rank](c, "/rank", token)
	}
	if err != nil {
		return nil, err
	}
	return &commonResponse.Data, nil
}

// NewStudent 创建一个新的学生账户。
func (c *SpiderClientImpl) NewStudent(username string, password string) (Student, error) {
	return NewStudentImpl(username, password, c)
}

func NewSpiderClientImpl(baseUrl string, client http.Client) *SpiderClientImpl {
	if client.Timeout == 0 {
		client.Timeout = 15 * time.Second
	}
	return &SpiderClientImpl{baseUrl: baseUrl, client: client}
}
