package account

import (
	"cached_proxy/utils"
	"log"
)

// Service 定义了账户服务的接口。
type Service interface {
	// GetAccountByAccountID 获取账户信息
	GetAccountByAccountID(accountID string) (Account, error)
	// GetAccountByToken 获取账户信息
	GetAccountByToken(token string) (Account, error)
	// Login 更新账户信息
	Login(username string, password string) (string, error)
	// LockAccount 锁定账户，锁定后账户返回的状态将会是锁定状态，如果需要恢复账户，需要重新登陆
	LockAccount(accountID string) error
}

type ServiceImpl struct {
	accountRepo repository
}

func NewServiceImpl(accountRepo repository) *ServiceImpl {
	return &ServiceImpl{accountRepo: accountRepo}
}

func (s *ServiceImpl) GetAccountByAccountID(accountID string) (Account, error) {
	account, err := s.accountRepo.GetAccountByAccountID(accountID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return account, nil
}

func (s *ServiceImpl) GetAccountByToken(token string) (Account, error) {
	account, err := s.accountRepo.GetAccountByToken(token)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return account, nil
}

// newAccount 创建一个新的账户
func (s *ServiceImpl) newAccount(username string, password string) (Account, error) {
	token, err := utils.GenerateUUID()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	account := &SimpleAccountImpl{
		Username:    username,
		Password:    password,
		StaticToken: token,
		status:      Normal,
	}
	return account, nil
}

func (s *ServiceImpl) Login(username string, password string) (string, error) {
	for {
		account, err := s.newAccount(username, password)
		if err != nil {
			return "", err
		}
		err = s.accountRepo.SaveOrUpdateAccount(account)
		if err != nil {
			if err.Error() == "StaticToken has been occupied" {
				continue
			}
			return "", err
		}
		return account.Token(), nil
	}
}

func (s *ServiceImpl) LockAccount(accountID string) error {
	account, err := s.GetAccountByAccountID(accountID)
	if err != nil {
		return err
	}
	account.setStatus(Banned)
	return s.accountRepo.SaveOrUpdateAccount(account)
}
