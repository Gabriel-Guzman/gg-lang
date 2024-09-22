package builtin

import (
	"fmt"
	"gg-lang/src/variables"
)

type Func interface {
	Name() string
	Call(args ...*variables.RuntimeValue) (*variables.RuntimeValue, error)
}

type Print struct{}

func (p *Print) Name() string {
	return "print"
}
func (p *Print) Call(args ...*variables.RuntimeValue) (*variables.RuntimeValue, error) {
	for _, arg := range args {
		fmt.Println(arg.Val)
	}
	return &variables.RuntimeValue{
		Val: nil,
		Typ: variables.Void,
	}, nil
}

func Defaults() []Func {
	return []Func{
		&Print{},
		&Length{},
	}
}
