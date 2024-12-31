package cache

// Updater 更新函数
type Updater[K, V any] interface {
	Invoke(key K) (V, error) // 更新执行入口
}

type LambdaUpdater[V any] struct {
	Invoker func(key string) (V, error)
}

func (l *LambdaUpdater[V]) Invoke(key string) (V, error) {
	return l.Invoker(key)
}
