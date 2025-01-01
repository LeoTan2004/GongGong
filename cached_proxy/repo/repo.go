package repo

// KVRepo 键值对存储接口
type KVRepo[K string, V any] interface {
	Get(key K) (value V, found bool)
	Set(key K, data V)
}

type MemRepo[K string, V any] struct {
	items map[K]V // 集合
}

func NewMemRepo[K string, V any]() *MemRepo[K, V] {
	return &MemRepo[K, V]{
		items: make(map[K]V),
	}
}

func (m *MemRepo[K, V]) Get(key K) (value V, found bool) {
	item, found := m.items[key]
	return item, found
}

func (m *MemRepo[K, V]) Set(key K, data V) {
	m.items[key] = data
}
