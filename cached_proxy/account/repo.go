package account

import (
	"cached_proxy/repo"
	"fmt"
)

type repository interface {
	// GetAccountByAccountID 获取账户信息
	GetAccountByAccountID(accountID string) (Account, error)
	// GetAccountByToken 获取账户信息
	GetAccountByToken(token string) (Account, error)
	// SaveOrUpdateAccount 保存或更新账户信息
	SaveOrUpdateAccount(account Account) error
}

type MemRepository struct {
	idRepo    repo.MemRepo[string, Account] // 用于根据账户ID查找账户
	tokenRepo repo.MemRepo[string, Account] // 用于根据token查找账户
}

func (m *MemRepository) GetAccountByAccountID(accountID string) (Account, error) {
	account, found := m.idRepo.Get(accountID)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (m *MemRepository) GetAccountByToken(token string) (Account, error) {
	account, found := m.tokenRepo.Get(token)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (m *MemRepository) SaveOrUpdateAccount(account Account) error {
	token := account.Token()
	accountId := account.AccountID()

	// if the token has been used by other account, we should reject the request
	tokenAccount, found := m.tokenRepo.Get(token)
	if found && tokenAccount.AccountID() != accountId {
		return fmt.Errorf("your token has been used by other account")
	}

	// if the account has other token, we should delete the old token
	formerAccount, found := m.idRepo.Get(accountId)
	if found && formerAccount.Token() != token {
		m.tokenRepo.Delete(formerAccount.Token())
	}

	// we will set the new account with the new token
	m.tokenRepo.Set(token, account)
	m.idRepo.Set(accountId, account)
	return nil
}
