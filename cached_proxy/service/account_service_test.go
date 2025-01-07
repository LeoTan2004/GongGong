package service

import (
	"cached_proxy/repo"
	"testing"
)

type MockAccount struct {
	id     string
	token  string
	statue AccountStatus
}

func (m *MockAccount) Status() AccountStatus {
	return m.statue
}

func (m *MockAccount) AccountId() string {
	return m.id
}

func (m *MockAccount) Token() string {
	return m.token
}

func TestGetAccountByAccountId(t *testing.T) {
	tests := []struct {
		name        string
		accountID   string
		setup       func(idRepo, tokenRepo repo.KVRepo[string, Account])
		expectError bool
	}{
		{
			name:      "Existing account",
			accountID: "user1",
			setup: func(idRepo, tokenRepo repo.KVRepo[string, Account]) {
				account := &MockAccount{id: "user1", token: "token1"}
				idRepo.Set("user1", account)
				tokenRepo.Set("token1", account)
			},
			expectError: false,
		},
		{
			name:        "Non-existing account",
			accountID:   "nonexistent",
			setup:       func(idRepo, tokenRepo repo.KVRepo[string, Account]) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idRepo := repo.NewMemRepo[string, Account]()
			tokenRepo := repo.NewMemRepo[string, Account]()
			service := &AccountServiceImpl{idRepo: idRepo, tokenRepo: tokenRepo}

			tt.setup(idRepo, tokenRepo)

			acc, err := service.GetAccountByAccountId(tt.accountID)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for accountID '%s', got nil", tt.accountID)
				}
			} else {
				if err != nil || acc.AccountId() != tt.accountID {
					t.Errorf("expected to find account with id '%s', got error: %v", tt.accountID, err)
				}
			}
		})
	}
}

func TestGetAccountByToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		setup       func(idRepo, tokenRepo repo.KVRepo[string, Account])
		expectError bool
	}{
		{
			name:  "Existing account",
			token: "token1",
			setup: func(idRepo, tokenRepo repo.KVRepo[string, Account]) {
				account := &MockAccount{id: "user1", token: "token1"}
				idRepo.Set("user1", account)
				tokenRepo.Set("token1", account)
			},
			expectError: false,
		},
		{
			name:        "Non-existing account",
			token:       "nonexistent",
			setup:       func(idRepo, tokenRepo repo.KVRepo[string, Account]) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idRepo := repo.NewMemRepo[string, Account]()
			tokenRepo := repo.NewMemRepo[string, Account]()
			service := &AccountServiceImpl{idRepo: idRepo, tokenRepo: tokenRepo}

			tt.setup(idRepo, tokenRepo)

			acc, err := service.GetAccountByToken(tt.token)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for token '%s', got nil", tt.token)
				}
			} else {
				if err != nil || acc.Token() != tt.token {
					t.Errorf("expected to find account with token '%s', got error: %v", tt.token, err)
				}
			}
		})
	}
}

func TestSaveOrUpdateAccount(t *testing.T) {
	tests := []struct {
		name        string
		account     *MockAccount
		setup       func(idRepo, tokenRepo repo.KVRepo[string, Account])
		expectError bool
	}{
		{
			name:        "Unique token",
			account:     &MockAccount{id: "user2", token: "token2"},
			setup:       func(idRepo, tokenRepo repo.KVRepo[string, Account]) {},
			expectError: false,
		},
		{
			name:    "Duplicate token",
			account: &MockAccount{id: "user3", token: "token1"},
			setup: func(idRepo, tokenRepo repo.KVRepo[string, Account]) {
				account := &MockAccount{id: "user1", token: "token1"}
				idRepo.Set("user1", account)
				tokenRepo.Set("token1", account)
			},
			expectError: true,
		},
		{
			name:    "Existing account with new token",
			account: &MockAccount{id: "user1", token: "token3"},
			setup: func(idRepo, tokenRepo repo.KVRepo[string, Account]) {
				account := &MockAccount{id: "user1", token: "token1"}
				idRepo.Set("user1", account)
				tokenRepo.Set("token1", account)
			},
			expectError: false,
		},
		{
			name:    "Existing account with same token",
			account: &MockAccount{id: "user1", token: "token3"},
			setup: func(idRepo, tokenRepo repo.KVRepo[string, Account]) {
				account := &MockAccount{id: "user1", token: "token3"}
				idRepo.Set("user1", account)
				tokenRepo.Set("token3", account)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idRepo := repo.NewMemRepo[string, Account]()
			tokenRepo := repo.NewMemRepo[string, Account]()
			service := &AccountServiceImpl{idRepo: idRepo, tokenRepo: tokenRepo}

			tt.setup(idRepo, tokenRepo)

			err := service.SaveOrUpdateAccount(tt.account)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for account '%v', got nil", tt.account)
				}
			} else {
				if err != nil {
					t.Errorf("expected to save/update account '%v', got error: %v", tt.account, err)
				}
			}
			if err == nil {
				acc, _ := idRepo.Get(tt.account.id)
				if acc.Token() != tt.account.token {
					t.Errorf("expected account with id '%s' to have token '%s', got '%s'", tt.account.id, tt.account.token, acc.Token())
				}
				account, _ := tokenRepo.Get(tt.account.token)
				if account.AccountId() != tt.account.id {
					t.Errorf("expected account with token '%s' to have id '%s', got '%s'", tt.account.token, tt.account.id, account.AccountId())
				}
			}
		})
	}
}
