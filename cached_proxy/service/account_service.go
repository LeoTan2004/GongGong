package service

import (
	"cached_proxy/repo"
	"cached_proxy/utils"
	"fmt"
)

type AccountStatus int

const (
	Normal = iota
	Banned
)

type Account interface {
	AccountId() string
	Token() string
	Status() AccountStatus
	setStatus(status AccountStatus)
}

type AccountService[A Account] interface {
	// GetAccountByAccountId 获取账户信息
	GetAccountByAccountId(accountID string) (*A, error)

	// GetAccountByToken 获取账户信息
	GetAccountByToken(token string) (*A, error)

	// SaveOrUpdateAccount 保存或更新账户信息
	SaveOrUpdateAccount(account A) error

	// LockAccount 锁定账户
	LockAccount(accountID string) error

	// GetUniqueToken 获取唯一的token
	GetUniqueToken() string
}

type AccountServiceImpl[A Account] struct {
	idRepo    repo.KVRepo[string, A]
	tokenRepo repo.KVRepo[string, A]
}

func (a *AccountServiceImpl[A]) GetUniqueToken() string {
	for {
		uuid, err := utils.GenerateUUID()
		if err != nil {
			return ""
		}
		_, found := a.tokenRepo.Get(uuid)
		if !found {
			return uuid
		}
	}

}

func (a *AccountServiceImpl[A]) LockAccount(accountID string) error {
	account, found := a.idRepo.Get(accountID)
	if !found {
		return fmt.Errorf("account not found")
	}
	account.setStatus(Banned)
	if err := a.SaveOrUpdateAccount(account); err != nil {
		return err
	}
	return nil
}

func (a *AccountServiceImpl[A]) GetAccountByAccountId(accountID string) (*A, error) {
	account, found := a.idRepo.Get(accountID)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return &account, nil
}

func (a *AccountServiceImpl[A]) GetAccountByToken(token string) (*A, error) {
	account, found := a.tokenRepo.Get(token)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return &account, nil
}

func (a *AccountServiceImpl[A]) SaveOrUpdateAccount(account A) error {
	token := account.Token()
	accountId := account.AccountId()

	// if the token has been used by other account, we should reject the request
	tokenAccount, found := a.tokenRepo.Get(token)
	if found && tokenAccount.AccountId() != accountId {
		return fmt.Errorf("your token has been used by other account")
	}

	// if the account has other token, we should delete the old token
	formerAccount, found := a.idRepo.Get(accountId)
	if found && formerAccount.Token() != token {
		a.tokenRepo.Delete(formerAccount.Token())
	}

	// we will set the new account with the new token
	a.tokenRepo.Set(token, account)
	a.idRepo.Set(accountId, account)
	return nil
}
