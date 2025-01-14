package cache

import "time"

// cacheItem 表示缓存的单个条目
type cacheItem[V any] struct {
	data     V         // 缓存的数据
	updateAt time.Time // 缓存更新时间
	submitAt time.Time // 提交更新时间
}
