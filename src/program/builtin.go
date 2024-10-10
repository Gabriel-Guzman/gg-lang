package program

import (
	"fmt"
	"gg-lang/src/variable"
)

type Func interface {
	Name() string
	Call(args ...*variable.RuntimeValue) (*variable.RuntimeValue, error)
}

type Print struct{}

func (p *Print) Name() string {
	return "print"
}
func (p *Print) Call(args ...*variable.RuntimeValue) (*variable.RuntimeValue, error) {
	for _, arg := range args {
		fmt.Println(arg.Val)
	}
	return &variable.RuntimeValue{
		Val: nil,
		Typ: variable.Void,
	}, nil
}

func Defaults() []Func {
	return []Func{
		&Print{},
		&Length{},
	}
}
