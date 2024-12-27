package cache

// ItemValidator  缓存校验器
type ItemValidator[V any] interface {
	Valid(item cacheItem[V]) bool // 校验缓存有效性
}
