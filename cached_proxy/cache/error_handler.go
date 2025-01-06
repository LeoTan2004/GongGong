package cache

type ErrorHandler[K any] interface {
	HandlerError(key K, err error)
}
