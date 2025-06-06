package gg

import (
	"errors"
	"fmt"
)

func Handle(err error) {
	if err == nil {
		return
	}
	var chillErr *RuntimeErr
	var critErr *CritErr
	switch {
	case errors.As(err, &chillErr):
		panic(fmt.Sprintf("Runtime error: %s\n", chillErr.Error()))
	case errors.As(err, &critErr):
		panic(fmt.Sprintf("Crit error: %s\n", critErr.Error()))
	default:
		panic(fmt.Sprintf("Unknown error (please wrap in ggErrs): %v\n", err))
	}
}
