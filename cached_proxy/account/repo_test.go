package account

import (
	"cached_proxy/repo"
	"testing"
)

func setup() *MemRepository {
	return &MemRepository{
		idRepo:    *repo.NewMemRepo[string, Account](),
		tokenRepo: *repo.NewMemRepo[string, Account](),
	}
}

func TestSaveOrUpdateAccount(t *testing.T) {
	memRepo := setup()

	account1 := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}

	err := memRepo.SaveOrUpdateAccount(account1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetAccountByAccountID(t *testing.T) {
	memRepo := setup()

	account1 := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}

	err := memRepo.SaveOrUpdateAccount(account1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	acc, err := memRepo.GetAccountByAccountID("user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.AccountID() != "user1" {
		t.Fatalf("expected accountID to be 'user1', got %v", acc.AccountID())
	}
}

func TestGetAccountByToken(t *testing.T) {
	memRepo := setup()

	account1 := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}

	err := memRepo.SaveOrUpdateAccount(account1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	acc, err := memRepo.GetAccountByToken("token1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.Token() != "token1" {
		t.Fatalf("expected token to be 'token1', got %v", acc.Token())
	}
}

func TestSaveOrUpdateAccountWithDuplicateToken(t *testing.T) {
	memRepo := setup()

	account1 := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}

	account2 := &simpleAccountImpl{
		accountID: "user2",
		token:     "token2",
		password:  "password2",
		status:    Normal,
	}

	err := memRepo.SaveOrUpdateAccount(account1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = memRepo.SaveOrUpdateAccount(account2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	account2.setToken("token1")
	err = memRepo.SaveOrUpdateAccount(account2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSaveOrUpdateAccountWithUpdatingToken(t *testing.T) {
	memRepo := setup()

	account1 := &simpleAccountImpl{
		accountID: "user1",
		token:     "token1",
		password:  "password1",
		status:    Normal,
	}

	err := memRepo.SaveOrUpdateAccount(account1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	account1.setToken("newToken1")
	err = memRepo.SaveOrUpdateAccount(account1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	acc, err := memRepo.GetAccountByToken("newToken1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.Token() != "newToken1" {
		t.Fatalf("expected token to be 'newToken1', got %v", acc.Token())
	}
}
