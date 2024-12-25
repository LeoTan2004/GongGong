package repo

// Updater 更新函数
type Updater[V any] interface {
	Invoke(params ...interface{}) (V, error) // 更新执行入口
}
