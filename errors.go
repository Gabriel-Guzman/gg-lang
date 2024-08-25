package main

import "fmt"

type runtimeError struct {
	msg string
}

func (err *runtimeError) Error() string {
	return fmt.Sprintf("runtime error: %s", err.msg)
}

type internalError struct {
	msg   string
	error error
}

func (err *internalError) Error() string {
	return fmt.Sprintf("internal error: %s (caused by: %v)", err.msg, err.error)
}
