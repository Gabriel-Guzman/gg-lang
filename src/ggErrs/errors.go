package ggErrs

import (
	"fmt"
	"runtime"
)

type errMessageFactory struct {
	text string
}

type SyntaxErr struct {
	Message string
}

func (err *SyntaxErr) Error() string {
	return err.Message
}

type RuntimeErr struct {
	Message string
}

func (err *RuntimeErr) Error() string {
	return err.Message
}

func Runtime(msg string, args ...interface{}) *RuntimeErr {
	return &RuntimeErr{Message: fmt.Sprintf(msg, args...)}
}

type CritErr struct {
	msg string
}

func (err *CritErr) Error() string {
	return err.msg
}

func Syntax(msg string, args ...interface{}) *SyntaxErr {
	return &SyntaxErr{Message: fmt.Sprintf(msg, args...)}
}

func Crit(msg string, args ...interface{}) *CritErr {
	pc, file, no, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return &CritErr{msg: fmt.Sprintf("%s:%d @ %s\n", file, no, details.Name()) + fmt.Sprintf(msg, args...)}
	}
	return &CritErr{msg: fmt.Sprintf(msg, args...)}
}

func WrappedCrit(msg string, args ...interface{}) *CritErr {
	pc, file, no, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return &CritErr{msg: fmt.Sprintf("%s:%d @ %s\n", file, no, details.Name()) + fmt.Sprintf(msg, args...)}
	}
	return &CritErr{msg: fmt.Sprintf(msg, args...)}
}
