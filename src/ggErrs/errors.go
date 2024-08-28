package ggErrs

import (
	"fmt"
	"runtime"
)

type ChillErr struct {
	Message string
}

func (err *ChillErr) Error() string {
	return fmt.Sprintf("ChillErr: %s", err.Message)
}

func Runtime(msg string, args ...interface{}) *ChillErr {
	return &ChillErr{Message: fmt.Sprintf(msg, args...)}
}

type CritErr struct {
	msg string
}

func (err *CritErr) Error() string {
	return fmt.Sprintf("CritErr: %s", err.msg)
}

func Crit(msg string, args ...interface{}) *CritErr {
	pc, file, no, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return &CritErr{msg: fmt.Sprintf("%s:%d @ %s\n", file, no, details.Name()) + fmt.Sprintf(msg, args...)}
	}
	return &CritErr{msg: fmt.Sprintf(msg, args...)}
}
