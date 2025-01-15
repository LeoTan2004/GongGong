package repo

import "sync"

// KVRepo 键值对存储接口
type KVRepo[K string, V any] interface {
	Get(key K) (value V, found bool)
	Set(key K, data V)
	Delete(key K) bool
}

type MemRepo[K string, V any] struct {
	items map[K]V      // 集合
	mu    sync.RWMutex // 读写锁
}

func (m *MemRepo[K, V]) Delete(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
	return true
}

func NewMemRepo[K string, V any]() *MemRepo[K, V] {
	return &MemRepo[K, V]{
		items: make(map[K]V),
	}
}

func (m *MemRepo[K, V]) Get(key K) (value V, found bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, found := m.items[key]
	return item, found
}

func (m *MemRepo[K, V]) Set(key K, data V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = data
}

type StaticRepo[K string, V any] struct {
	value V
}

func NewStaticRepo[K string, V any]() *StaticRepo[K, V] {
	return &StaticRepo[K, V]{}
}

func (s *StaticRepo[K, V]) Get(_ K) (value V, found bool) {
	return s.value, true
}

func (s *StaticRepo[K, V]) Set(_ K, data V) {
	s.value = data
}

func (s *StaticRepo[K, V]) Delete(_ K) bool {
	return true
}
