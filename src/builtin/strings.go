package builtin

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/variables"
)

/* type Func interface {
	Name() string
	Call(args ...*variables.RuntimeValue) *variables.RuntimeValue
} */

type Length struct{}

func (l *Length) Name() string {
	return "len"
}

func (l *Length) Call(args ...*variables.RuntimeValue) (*variables.RuntimeValue, error) {
	if len(args) != 1 {
		return nil, ggErrs.Runtime("len expects one argument")
	}

	switch args[0].Typ {
	case variables.String:
		return &variables.RuntimeValue{
			Val: len(args[0].Val.(string)),
			Typ: variables.Integer,
		}, nil
	default:
		return nil, ggErrs.Runtime("len argument must be a string, got %s", args[0].Typ.String())
	}
}
