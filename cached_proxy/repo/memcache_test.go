package repo

import (
	exec "cached_proxy/executor"
	"testing"
	"time"
)

// Mock ItemValidator
type mockValidator[V any] struct{}

func (m *mockValidator[V]) Valid(cacheItem[V]) bool {
	// 简单实现：假设所有缓存项都有效
	return false
}

// Mock Updater
type mockUpdater[V any] struct {
	data V
}

func (u *mockUpdater[V]) Invoke(...interface{}) (V, error) {
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
func TestMemCache_Get(t *testing.T) {
	repo := NewMockRepo[string, cacheItem[string]]()
	validator := &mockValidator[string]{}
	updater := &mockUpdater[string]{data: "updated value"}
	executor := exec.NewWorkerPool(4)
	executor.Run()
	cache := NewMemCache(validator, updater, repo, executor)

	// Set an initial value in the cache
	cache.Set("key1", "initial value")

	// Test retrieving the value
	data, found := cache.Get("key1")
	if !found || data != "initial value" {
		t.Errorf("Get() failed, expected 'initial value', got '%v' (found: %v)", data, found)
	}

	// Simulate expired value
	time.Sleep(1 * time.Second) // Let the item become "stale"
	data, found = cache.Get("key1")
	if !found || data != "updated value" {
		t.Errorf("Get() failed, expected 'updated value', got '%v' (found: %v)", data, found)
	}
}

// Test for Set
func TestMemCache_Set(t *testing.T) {
	repo := NewMockRepo[string, cacheItem[string]]()
	validator := &mockValidator[string]{}
	updater := &mockUpdater[string]{data: "updated value"}
	executor := exec.NewWorkerPool(4)
	cache := NewMemCache(validator, updater, repo, executor)

	// Add a new value
	cache.Set("key1", "test value")

	// Verify the value is stored correctly
	item, found := repo.Get("key1")
	if !found || item.data != "test value" {
		t.Errorf("Set() failed, expected 'test value', got '%v' (found: %v)", item.data, found)
	}
}
