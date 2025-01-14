package account

import (
	"errors"
	"testing"
)

// MockRepository is a mock implementation of the repository interface
type MockRepository struct {
	accountsByID    map[string]Account
	accountsByToken map[string]Account
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		accountsByID:    make(map[string]Account),
		accountsByToken: make(map[string]Account),
	}
}

func (m *MockRepository) GetAccountByAccountID(accountID string) (Account, error) {
	account, found := m.accountsByID[accountID]
	if !found {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (m *MockRepository) GetAccountByToken(token string) (Account, error) {
	account, found := m.accountsByToken[token]
	if !found {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (m *MockRepository) SaveOrUpdateAccount(account Account) error {
	m.accountsByID[account.AccountID()] = account
	m.accountsByToken[account.Token()] = account
	return nil
}

func TestServiceImpl_GetAccountByAccountID(t *testing.T) {
	mockRepo := NewMockRepository()
	service := &ServiceImpl{accountRepo: mockRepo}

	account := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}
	err := mockRepo.SaveOrUpdateAccount(account)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	acc, err := service.GetAccountByAccountID("user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.AccountID() != "user1" {
		t.Fatalf("expected accountID to be 'user1', got %v", acc.AccountID())
	}
}

func TestServiceImpl_GetAccountByToken(t *testing.T) {
	mockRepo := NewMockRepository()
	service := &ServiceImpl{accountRepo: mockRepo}

	account := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}
	err := mockRepo.SaveOrUpdateAccount(account)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	acc, err := service.GetAccountByToken("token1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.Token() != "token1" {
		t.Fatalf("expected token to be 'token1', got %v", acc.Token())
	}
}

func TestServiceImpl_Login(t *testing.T) {
	mockRepo := NewMockRepository()
	service := &ServiceImpl{accountRepo: mockRepo}

	account := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}
	err := mockRepo.SaveOrUpdateAccount(account)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := service.Login("user1", "password1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "token1" {
		t.Fatalf("expected new token except %v", token)
	}
}

func TestServiceImpl_LockAccount(t *testing.T) {
	mockRepo := NewMockRepository()
	service := &ServiceImpl{accountRepo: mockRepo}

	account := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}
	err := mockRepo.SaveOrUpdateAccount(account)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = service.LockAccount("user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	acc, err := service.GetAccountByAccountID("user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.Status() != Banned {
		t.Fatalf("expected status to be 'Banned', got %v", acc.Status())
	}
}
