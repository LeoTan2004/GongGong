package cache

import (
	"cached_proxy/err_handler"
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
func NewReadOnlyCache[K string, V any](valid ItemValidator[V], updater Updater[K, V], repo repo.KVRepo[K, cacheItem[V]], executor executor.Executor) *ReadOnlyCache[K, V] {
	return &ReadOnlyCache[K, V]{
		items:    repo,    // 初始化缓存映射
		valid:    valid,   // 校验器
		updater:  updater, // 更新器
		executor: executor,
	}
}

// ReadOnlyCache 定义缓存结构
type ReadOnlyCache[K string, V any] struct {
	items         repo.KVRepo[K, cacheItem[V]] // 缓存项的存储映射
	valid         ItemValidator[V]             // 缓存有效性校验器
	updater       Updater[K, V]                // 缓存更新器
	executor      executor.Executor            // 执行器
	onUpdateError err_handler.ErrorHandler     // 更新错误处理器
}

// Get 获取指定键的缓存数据
//
// 返回值：数据、是否存在、是否过期
func (c *ReadOnlyCache[K, V]) Get(key K) (data V, found bool) {
	item, found, needUpdate := c.getWithValid(key)
	if !needUpdate {
		// updateTask 更新任务
		item.submitAt = time.Now()
		updateTask := c.getUpdaterTask(key)
		c.executor.Submit(updateTask)
	}
	return item.data, found
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
			c.onUpdateError.HandlerError(err)
		} else {
			c.Set(key, result)
		}
	}
	return updateTask
}

// 获取数据并且检查有效性
func (c *ReadOnlyCache[string, V]) getWithValid(key string) (item cacheItem[V], found bool, needUpdate bool) {
	item, found = c.items.Get(key) // 在缓存中查找 key
	if !found {
		return cacheItem[V]{}, found, true
	}
	return item, found, c.valid.Valid(item)
}
