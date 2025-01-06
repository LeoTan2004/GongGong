package cache

import "time"

type ItemStatus int

const (
	Valid    = iota // 缓存有效
	Expired         // 缓存过期
	NotFound        // 缓存未找到
	Updating        // 缓存正在更新
)

// StatusChecker  缓存校验器
type StatusChecker[V any] interface {
	StatusOf(item *cacheItem[V]) ItemStatus // 校验缓存状态
}

type DailyStatusChecker[V any] struct {
	submitExpireInterval time.Duration // 提交过期时间间隔
}

func (d *DailyStatusChecker[V]) StatusOf(item *cacheItem[V]) ItemStatus {
	// 判断是否为同一天的数据
	if item == nil {
		return NotFound
	}
	if item.updateAt.Day() == time.Now().Day() {
		return Valid
	}
	if item.submitAt.Add(d.submitExpireInterval).After(time.Now()) {
		return Updating
	}
	return Expired

}

// NewDailyStatusChecker 创建默认缓存校验器
func NewDailyStatusChecker[V any](submitExpireAt time.Duration) *DailyStatusChecker[V] {
	return &DailyStatusChecker[V]{
		submitExpireInterval: submitExpireAt,
	}
}

// IntervalStatusChecker 基于时间间隔的缓存校验器
type IntervalStatusChecker[V any] struct {
	updateExpireInterval time.Duration // 更新过期时间间隔
	submitExpireInterval time.Duration // 提交过期时间间隔
}

func (d *IntervalStatusChecker[V]) StatusOf(item *cacheItem[V]) ItemStatus {
	if item == nil {
		return NotFound
	}
	if item.updateAt.Add(d.updateExpireInterval).After(time.Now()) {
		return Valid
	}
	if item.submitAt.Add(d.submitExpireInterval).After(time.Now()) {
		return Updating
	}
	return Expired

}

// NewIntervalStatusChecker 创建默认缓存校验器
func NewIntervalStatusChecker[V any](updateExpireAt, submitExpireAt time.Duration) *IntervalStatusChecker[V] {
	return &IntervalStatusChecker[V]{
		updateExpireInterval: updateExpireAt,
		submitExpireInterval: submitExpireAt,
	}
}
