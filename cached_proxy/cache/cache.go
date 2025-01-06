package cache

import (
	"cached_proxy/executor"
	"cached_proxy/repo"
	"time"
)

// cacheItem 表示缓存的单个条目
type cacheItem[V any] struct {
	data     V         // 缓存的数据
	updateAt time.Time // 缓存更新时间
	submitAt time.Time // 提交更新时间
}

// NewReadOnlyCache 创建并初始化一个新的缓存实例
func NewReadOnlyCache[K string, V any](statusChecker StatusChecker[V], updater Updater[K, V], executor executor.Executor, handler ErrorHandler[K]) *ReadOnlyCache[K, V] {
	return &ReadOnlyCache[K, V]{
		items:         repo.NewMemRepo[K, cacheItem[V]](), // 初始化缓存映射
		statusChecker: statusChecker,                      // 检查器
		updater:       updater,                            // 更新器
		executor:      executor,                           // 执行器
		onUpdateError: handler,                            // 错误处理器
	}
}

// ReadOnlyCache 定义缓存结构
type ReadOnlyCache[K string, V any] struct {
	items         repo.KVRepo[K, cacheItem[V]] // 缓存项的存储映射
	statusChecker StatusChecker[V]             // 缓存状态检查器
	updater       Updater[K, V]                // 缓存更新器
	executor      executor.Executor            // 执行器
	onUpdateError ErrorHandler[K]              // 更新错误处理器
}

// Get 获取指定键的缓存数据
//
// 返回值：数据、是否最新、是否过期
func (c *ReadOnlyCache[K, V]) Get(key K) (value V, valid bool) {
	item, status := c.getWithValid(key)
	needUpdate := status == Expired || status == NotFound
	if needUpdate {
		item.submitAt = time.Now()
		c.executor.Submit(c.getUpdaterTask(key))
	}
	return item.data, status == Valid
}

// Set 将数据添加到缓存或更新现有数据
//
// 返回值：键、数据
func (c *ReadOnlyCache[K, V]) Set(key K, data V) {
	item, found := c.items.Get(key) // 在缓存中查找 key
	if !found {
		item = cacheItem[V]{
			data:     data,
			updateAt: time.Now(), // 当前时间加上 TTL
			submitAt: time.Now(),
		}

	} else {
		item = cacheItem[V]{
			data:     data,
			updateAt: time.Now(), // 当前时间加上 TTL
			submitAt: item.submitAt,
		}
	}
	c.items.Set(key, item)
}

// 获取更新任务，更新任务中会调用创建时给的更新器去更新。同时也会自动处理好时间记录等问题
func (c *ReadOnlyCache[string, V]) getUpdaterTask(key string) func() {
	updateTask := func() {
		result, err := c.updater.Invoke(key)
		if err != nil {
			c.onUpdateError.HandlerError(key, err)
		} else {
			c.Set(key, result)
		}
	}
	return updateTask
}

// 获取数据并且检查有效性
func (c *ReadOnlyCache[string, V]) getWithValid(key string) (item cacheItem[V], status ItemStatus) {
	item, _ = c.items.Get(key) // 在缓存中查找 key
	return item, c.statusChecker.StatusOf(&item)
}
