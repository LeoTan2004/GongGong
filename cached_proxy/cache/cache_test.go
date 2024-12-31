package cache

import (
	"cached_proxy/executor"
	"testing"
	"time"
)

// Mock ItemValidator
type mockValidator[V any] struct{}

func (m *mockValidator[V]) Valid(*cacheItem[V]) bool {
	// 简单实现：假设所有缓存项都有效
	return false
}

// Mock Updater
type mockUpdater[K, V any] struct {
	data V
}

func (u *mockUpdater[K, V]) Invoke(K) (V, error) {
	return u.data, nil
}

// Mock KVRepo
type MockRepo[K comparable, V any] struct {
	items map[K]V
}

func NewMockRepo[K comparable, V any]() *MockRepo[K, V] {
	return &MockRepo[K, V]{
		items: make(map[K]V),
	}
}

func (r *MockRepo[K, V]) Get(key K) (V, bool) {
	val, ok := r.items[key]
	return val, ok
}

func (r *MockRepo[K, V]) Set(key K, value V) {
	r.items[key] = value
}

// Test for Get
func TestCache_Get(t *testing.T) {
	repo := NewMockRepo[string, cacheItem[string]]()
	validator := &mockValidator[string]{}
	updater := &mockUpdater[string, string]{data: "updated value"}
	exec := executor.NewWorkerPool(4)
	exec.Run()
	cache := NewReadOnlyCache(validator, updater, repo, exec)

	// Set an initial value in the cache
	cache.Set("key1", "initial value")

	// Test retrieving the value
	data, valid := cache.Get("key1")
	if !valid || data != "initial value" {
		t.Errorf("Get() failed, expected 'initial value', got '%v' (valid: %v)", data, valid)
	}

	// Simulate expired value
	time.Sleep(10 * time.Microsecond) // Let the item become "stale"
	data, valid = cache.Get("key1")
	if !valid || data != "updated value" {
		t.Errorf("Get() failed, expected 'updated value', got '%v' (valid: %v)", data, valid)
	}
}

// Test for Set
func TestCache_Set(t *testing.T) {
	repo := NewMockRepo[string, cacheItem[string]]()
	validator := &mockValidator[string]{}
	updater := &mockUpdater[string, string]{data: "updated value"}
	exec := executor.NewWorkerPool(4)
	cache := NewReadOnlyCache(validator, updater, repo, exec)

	// Add a new value
	cache.Set("key1", "test value")

	// Verify the value is stored correctly
	item, valid := repo.Get("key1")
	if !valid || item.data != "test value" {
		t.Errorf("Set() failed, expected 'test value', got '%v' (valid: %v)", item.data, valid)
	}
}
