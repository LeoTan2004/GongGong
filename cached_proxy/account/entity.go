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
}

type simpleAccountImpl struct {
	accountID string
	token     string
	password  string
	status    Status
}

func (s *simpleAccountImpl) AccountID() string {
	return s.accountID
}

func (s *simpleAccountImpl) Token() string {
	return s.token
}

func (s *simpleAccountImpl) Status() Status {
	return s.status
}

func (s *simpleAccountImpl) setStatus(status Status) {
	s.status = status
}

func (s *simpleAccountImpl) setToken(token string) {
	s.token = token
}
