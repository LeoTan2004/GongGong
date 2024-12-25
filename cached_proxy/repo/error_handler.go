package repo

type ErrorHandler interface {
	HandlerError(err error)
}
