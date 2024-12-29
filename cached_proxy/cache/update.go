package cache

// Updater 更新函数
type Updater[K, V any] interface {
	Invoke(key K) (V, error) // 更新执行入口
}
