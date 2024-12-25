package repo

import (
	exec "cached_proxy/executor"
	"time"
)

// cacheItem 表示缓存的单个条目
type cacheItem[V any] struct {
	data     V         // 缓存的数据
	updateAt time.Time // 缓存更新时间
	submitAt time.Time // 提交更新时间
}

// NewMemCache 创建并初始化一个新的缓存实例
func NewMemCache[K string, V any](valid ItemValidator[V], updater Updater[V], repo KVRepo[K, cacheItem[V]], executor exec.Executor) *MemCache[K, V] {
	if repo == nil {
		repo = NewMemRepo[K, cacheItem[V]]()
	}
	return &MemCache[K, V]{
		items:    repo,    // 初始化缓存映射
		valid:    valid,   // 校验器
		updater:  updater, // 更新器
		executor: executor,
	}
}

// MemCache 定义缓存结构
type MemCache[K string, V any] struct {
	items         KVRepo[K, cacheItem[V]] // 缓存项的存储映射
	valid         ItemValidator[V]        // 缓存有效性校验器
	updater       Updater[V]              // 缓存更新器
	executor      exec.Executor           // 执行器
	onUpdateError ErrorHandler            // 更新错误处理器
}

// Get 获取指定键的缓存数据
//
// 返回值：数据、是否存在、是否过期
func (c *MemCache[K, V]) Get(key K) (data V, found bool) {
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
func (c *MemCache[K, V]) Set(key K, data V) {
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
func (c *MemCache[string, V]) getUpdaterTask(key string) func() {
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
func (c *MemCache[string, V]) getWithValid(key string) (item cacheItem[V], found bool, needUpdate bool) {
	item, found = c.items.Get(key) // 在缓存中查找 key
	if !found {
		return cacheItem[V]{}, found, true
	}
	return item, found, c.valid.Valid(item)
}
