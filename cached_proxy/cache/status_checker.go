package cache

import "time"

type ItemStatus int

const (
	Valid    = iota // 缓存有效
	Expired         // 缓存过期
	NotFound        // 缓存未找到
	Updating        // 缓存正在更新
)

// ItemStatusChecker  缓存校验器
type ItemStatusChecker[V any] interface {
	StatusOf(item *cacheItem[V]) ItemStatus // 校验缓存状态
}

type DefaultItemStatusChecker[V any] struct {
	updateExpireAt time.Duration // 更新过期时间
	submitExpireAt time.Duration // 提交过期时间
}

func (d *DefaultItemStatusChecker[V]) StatusOf(item *cacheItem[V]) ItemStatus {
	if item == nil {
		return NotFound
	}
	if item.updateAt.Add(d.updateExpireAt).After(time.Now()) {
		return Valid
	}
	if item.submitAt.Add(d.submitExpireAt).After(time.Now()) {
		return Updating
	}
	return Expired

}

// NewDefaultItemStatusChecker 创建默认缓存校验器
func NewDefaultItemStatusChecker[V any](updateExpireAt, submitExpireAt time.Duration) *DefaultItemStatusChecker[V] {
	return &DefaultItemStatusChecker[V]{
		updateExpireAt: updateExpireAt,
		submitExpireAt: submitExpireAt,
	}
}
