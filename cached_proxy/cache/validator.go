package cache

import "time"

// ItemValidator  缓存校验器
type ItemValidator[V any] interface {
	Valid(item *cacheItem[V]) bool // 校验缓存有效性
}

type DefaultItemValidator[V any] struct {
	updateExpireAt time.Duration // 更新过期时间
	submitExpireAt time.Duration // 提交过期时间
}

func (d *DefaultItemValidator[V]) Valid(item *cacheItem[V]) bool {
	if item == nil { // 如果缓存项为空，则返回 false
		return false
	}
	if item.updateAt.Add(d.updateExpireAt).After(time.Now()) { // 如果更新时间加上保质期大于当前时间，则返回true，表示缓存没过期，可以使用
		return true
	} else {
		// 如果提交时间加上保质期小于当前时间，则返回true，表示缓存没过期，可以使用，否则返回false，表示缓存过期
		return item.submitAt.Add(d.submitExpireAt).After(time.Now())
	}

}

// NewDefaultItemValidator 创建默认缓存校验器
func NewDefaultItemValidator[V any](updateExpireAt, submitExpireAt time.Duration) *DefaultItemValidator[V] {
	return &DefaultItemValidator[V]{
		updateExpireAt: updateExpireAt,
		submitExpireAt: submitExpireAt,
	}
}
