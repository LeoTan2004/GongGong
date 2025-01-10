package service

import (
	"cached_proxy/cache"
	"cached_proxy/executor"
	"cached_proxy/feign"
	"log"
)

type InformationService[V any] interface {
	// GetInfo 获取学生信息
	GetInfo(key string) (*V, bool)
}

type InformationServiceImpl[V any] struct {
	// 当处理器遇到错误时调用
	onAuthorityError func(username string)
	// 获取学生代理
	getStudentProxy func(username string) (*feign.Student, error)
	// 缓存校验结构
	checker cache.StatusChecker[V]
	// 缓存存储
	cache cache.ReadOnlyCache[string, V]
	// updater
	updater func(student *feign.Student) (*V, error)
}

type errorHandlerWrapper struct {
	onAuthorityError func(username string)
}

// 更新失败的错误处理
func (e *errorHandlerWrapper) errorHandler(key string, err error) {
	if err.Error() == "unauthorized" {
		log.Printf("account %s has no authority", key)
		e.onAuthorityError(key)
		return
	}
	log.Printf("failed to update %s: %v", key, err)
}

type updaterWrapper[V any] struct {
	service *InformationServiceImpl[V]
}

// cache缓存结构触发更新操作
func (u *updaterWrapper[V]) updateItem(key string) (*V, error) {
	p := u.service
	student, err := p.getStudentProxy(key)
	if err != nil {
		return nil, err
	}
	updater, err := p.updater(student)
	if err != nil {
		return nil, err
	}
	return updater, nil
}

func (p *InformationServiceImpl[V]) GetInfo(key string) (*V, bool) {
	value, valid := p.cache.Get(key)
	return &value, valid
}

func NewPublicInformationServiceImpl[V any](
	onAuthorityError func(username string), // 当处理器遇到错误时调用
	getStudentProxy func(username string) (*feign.Student, error), // 获取学生代理
	checker cache.StatusChecker[V], // 缓存校验
	updater func(student *feign.Student) (*V, error), // 更新器
	exec executor.Executor,
) *InformationServiceImpl[V] {
	ser := &InformationServiceImpl[V]{onAuthorityError: onAuthorityError, getStudentProxy: getStudentProxy, checker: checker, updater: updater}

	ser.cache = *cache.NewReadOnlyCache[string, V](
		ser.checker,
		(&updaterWrapper[V]{service: ser}).updateItem,
		exec,
		(&errorHandlerWrapper{onAuthorityError: onAuthorityError}).errorHandler,
	)
	return ser
}
