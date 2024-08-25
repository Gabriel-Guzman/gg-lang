package errors

import "fmt"

type RuntimeError struct {
	Message string
}

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error: %s", err.Message)
}

type InternalError struct {
	msg   string
	error error
}

func (err *InternalError) Error() string {
	return fmt.Sprintf("internal error: %s (caused by: %v)", err.msg, err.error)
}
