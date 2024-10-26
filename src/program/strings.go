package program

import (
	"gg-lang/src/gg"
	"gg-lang/src/variable"
)

/* type Func interface {
	Name() string
	Call(args ...*variables.RuntimeValue) *variables.RuntimeValue
} */

type Length struct{}

func (l *Length) Name() string {
	return "len"
}

func (l *Length) Call(args ...*variable.RuntimeValue) (*variable.RuntimeValue, error) {
	if len(args) != 1 {
		return nil, gg.Runtime("len expects one argument")
	}

	switch args[0].Typ {
	case variable.String:
		return &variable.RuntimeValue{
			Val: len(args[0].Val.(string)),
			Typ: variable.Integer,
		}, nil
	default:
		return nil, gg.Runtime("len argument must be a string, got %s", args[0].Typ.String())
	}
}
