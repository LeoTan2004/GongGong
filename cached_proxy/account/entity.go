package account

// Status 定义了账户状态的类型。
type Status int

const (
	Normal = iota
	Banned
)

// Account 定义了账户的接口。
type Account interface {
	// AccountID 获取账户的唯一标识符。
	AccountID() string
	// Token 获取账户的身份验证令牌。
	Token() string
	// setToken 设置账户的身份验证令牌。
	setToken(token string)
	// Status 获取账户的状态。
	Status() Status
	// setStatus 设置账户的状态。
	setStatus(status Status)
	// GetPassword 获取账户的密码。
	GetPassword() string
}

type SimpleAccountImpl struct {
	Username    string
	StaticToken string
	Password    string
	status      Status
}

func (s *SimpleAccountImpl) GetPassword() string {
	return s.Password
}

func (s *SimpleAccountImpl) AccountID() string {
	return s.Username
}

func (s *SimpleAccountImpl) Token() string {
	return s.StaticToken
}

func (s *SimpleAccountImpl) Status() Status {
	return s.status
}

func (s *SimpleAccountImpl) setStatus(status Status) {
	s.status = status
}

func (s *SimpleAccountImpl) setToken(token string) {
	s.StaticToken = token
}
