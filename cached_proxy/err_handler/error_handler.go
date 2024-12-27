package err_handler

type ErrorHandler interface {
	HandlerError(err error)
}
