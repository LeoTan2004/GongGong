package cache

import (
	"cached_proxy/executor"
	"fmt"
	"testing"
	"time"
)

// Mock StatusChecker
type mockChecker[V any] struct {
	status ItemStatus
}

func (m *mockChecker[V]) StatusOf(*cacheItem[V]) ItemStatus {
	// 简单实现：假设所有缓存项都有效
	return m.status
}

// Mock Updater
type mockUpdater[K, V any] struct {
	data      V
	returnErr bool
	Error     error
}

func (u *mockUpdater[K, V]) Invoke(K) (V, error) {
	if u.returnErr {
		return u.data, u.Error
	}
	return u.data, nil
}

type errorHandler[K any] struct {
	calledCount int
}

func (e *errorHandler[K]) HandlerError(_ K, _ error) {
	e.calledCount++
}

// Test for Get
func TestCache_Get(t *testing.T) {
	tests := []struct {
		name           string                  // 测试名称
		key            string                  // 键
		checker        StatusChecker[string]   // 状态检查器
		updater        Updater[string, string] // 更新器
		expected       string                  // 预期值
		again          bool                    // 是否再次获取
		expectedAgain  string                  // 再次获取的预期值
		expectedErrors int                     // 预期错误次数
	}{
		{
			name:     "Cache occurred and valid",
			key:      "key1",
			checker:  &mockChecker[string]{status: Valid},
			updater:  nil,
			expected: "initial value",
			again:    false,
		},
		{
			name:          "Cache occurred and expired, update success",
			key:           "key1",
			checker:       &mockChecker[string]{status: Expired},
			updater:       &mockUpdater[string, string]{data: "updated value"},
			expected:      "initial value",
			again:         true,
			expectedAgain: "updated value",
		},
		{
			name:    "Cache occurred and expired, update failed",
			key:     "key1",
			checker: &mockChecker[string]{status: Expired},
			updater: &mockUpdater[string, string]{
				data:      "updated value",
				returnErr: true,
				Error:     fmt.Errorf("mock error"),
			},
			expected:       "initial value",
			again:          true,
			expectedAgain:  "initial value",
			expectedErrors: 1,
		},
		{
			name:          "Cache not occurred and update success",
			key:           "not_exist_key",
			checker:       &mockChecker[string]{status: NotFound},
			updater:       &mockUpdater[string, string]{data: "updated value"},
			expected:      "",
			again:         true,
			expectedAgain: "updated value",
		},
		{
			name:    "Cache not occurred and update failed",
			key:     "not_exist_key",
			checker: &mockChecker[string]{status: NotFound},
			updater: &mockUpdater[string, string]{
				data:      "updated value",
				returnErr: true,
				Error:     fmt.Errorf("mock error"),
			},
			expected:       "",
			again:          true,
			expectedAgain:  "",
			expectedErrors: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := executor.NewWorkerPool(4)
			exec.Run()
			defer exec.Wait()
			onUpdaterError := &errorHandler[string]{}
			cache := NewReadOnlyCache(tt.checker, tt.updater, exec, onUpdaterError)
			repo := cache.items

			// Add a new value
			repo.Set("key1", cacheItem[string]{
				data:     "initial value",
				updateAt: time.Now().Add(-10 * time.Second),
				submitAt: time.Now().Add(-10 * time.Second),
			})

			// Get the value
			value, valid := cache.Get(tt.key)
			if value != tt.expected {
				t.Errorf("Get() failed, expected '%v', got '%v' (valid: %v)", tt.expected, value, valid)
			}
			time.Sleep(1 * time.Millisecond)
			// Get the value again
			if tt.again {
				value, valid = cache.Get(tt.key)
				if value != tt.expectedAgain {
					t.Errorf("Get() failed, expected '%v', got '%v' (valid: %v)", tt.expectedAgain, value, valid)
				}
			}

			// Check the error count
			if tt.expectedErrors != 0 {
				if onUpdaterError.calledCount != tt.expectedErrors {
					t.Errorf("Get() failed, expected %v errors, got %v", tt.expectedErrors, onUpdaterError.calledCount)
				}
			}
		})
	}
}

// Test for Set
func TestCache_Set(t *testing.T) {
	validator := &mockChecker[string]{}
	updater := &mockUpdater[string, string]{data: "updated value"}
	exec := executor.NewWorkerPool(4)
	cache := NewReadOnlyCache(validator, updater, exec, nil)
	repo := cache.items

	// Add a new value
	cache.Set("key1", "test value")

	// Verify the value is stored correctly
	item, valid := repo.Get("key1")
	if !valid || item.data != "test value" {
		t.Errorf("Set() failed, expected 'test value', got '%v' (valid: %v)", item.data, valid)
	}
}
