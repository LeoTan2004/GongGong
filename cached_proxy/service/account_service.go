package service

import (
	"cached_proxy/repo"
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
}

type AccountService[A Account] interface {
	// GetAccountByAccountId 获取账户信息
	GetAccountByAccountId(accountID string) (A, error)

	// GetAccountByToken 获取账户信息
	GetAccountByToken(token string) (A, error)

	// SaveOrUpdateAccount 保存或更新账户信息
	SaveOrUpdateAccount(account A) error
}

type AccountServiceImpl struct {
	idRepo    repo.KVRepo[string, Account]
	tokenRepo repo.KVRepo[string, Account]
}

func (a *AccountServiceImpl) GetAccountByAccountId(accountID string) (Account, error) {
	account, found := a.idRepo.Get(accountID)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (a *AccountServiceImpl) GetAccountByToken(token string) (Account, error) {
	account, found := a.tokenRepo.Get(token)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (a *AccountServiceImpl) SaveOrUpdateAccount(account Account) error {
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
