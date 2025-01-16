package account

import (
	"cached_proxy/repo"
	"encoding/gob"
	"fmt"
	path2 "path"
)

type repository interface {
	// GetAccountByAccountID 获取账户信息
	GetAccountByAccountID(accountID string) (Account, error)
	// GetAccountByToken 获取账户信息
	GetAccountByToken(token string) (Account, error)
	// SaveOrUpdateAccount 保存或更新账户信息
	SaveOrUpdateAccount(account Account) error
}

type Repository struct {
	idRepo    repo.KVRepo[string, Account] // 用于根据账户ID查找账户
	tokenRepo repo.KVRepo[string, Account] // 用于根据token查找账户
}

func NewMemRepository() *Repository {
	return &Repository{
		idRepo:    repo.NewMemRepo[string, Account](),
		tokenRepo: repo.NewMemRepo[string, Account](),
	}
}

func init() {
	gob.Register(&SimpleAccountImpl{})
}

func NewFileRepository(path string) *Repository {
	return &Repository{
		idRepo:    repo.NewFileRepos[string, Account](path2.Join(path, "account_id.gob")),
		tokenRepo: repo.NewFileRepos[string, Account](path2.Join(path, "account_token.gob")),
	}
}

func (m *Repository) GetAccountByAccountID(accountID string) (Account, error) {
	account, found := m.idRepo.Get(accountID)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (m *Repository) GetAccountByToken(token string) (Account, error) {
	account, found := m.tokenRepo.Get(token)
	if !found {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (m *Repository) SaveOrUpdateAccount(account Account) error {
	token := account.Token()
	accountId := account.AccountID()

	// if the StaticToken has been used by other account, we should reject the request
	tokenAccount, found := m.tokenRepo.Get(token)
	if found && tokenAccount.AccountID() != accountId {
		return fmt.Errorf("StaticToken has been occupied")
	}

	// if the account has other StaticToken, we should delete the old StaticToken
	formerAccount, found := m.idRepo.Get(accountId)
	if found && formerAccount.Token() != token {
		m.tokenRepo.Delete(formerAccount.Token())
	}

	// we will set the new account with the new StaticToken
	m.tokenRepo.Set(token, account)
	m.idRepo.Set(accountId, account)
	return nil
}
