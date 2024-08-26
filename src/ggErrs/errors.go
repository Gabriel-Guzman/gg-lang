package ggErrs

import "fmt"

type Runtime struct {
	Message string
}

func (err *Runtime) Error() string {
	return fmt.Sprintf("runtime error: %s", err.Message)
}

func NewRuntime(msg string, args ...interface{}) *Runtime {
	return &Runtime{Message: fmt.Sprintf(msg, args...)}
}

type Internal struct {
	msg   string
	error error
}

func (err *Internal) Error() string {
	switch err.error {
	case nil:
		return fmt.Sprintf("internal error: %s", err.msg)
	default:
		return fmt.Sprintf("internal error: %s (caused by non-gg code: %v)", err.msg, err.error)
	}
}

func NewInternal(err error, msg string, args ...interface{}) *Internal {
	return &Internal{msg: fmt.Sprintf(msg, args...), error: err}
}
