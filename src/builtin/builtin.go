package builtin

import (
	"fmt"
	"gg-lang/src/variables"
)

type Func interface {
	Name() string
	Call(args ...*variables.RuntimeValue) *variables.RuntimeValue
}

type Print struct{}

func (p *Print) Name() string {
	return "print"
}
func (p *Print) Call(args ...*variables.RuntimeValue) *variables.RuntimeValue {
	for _, arg := range args {
		fmt.Println(arg.Val)
	}
	return nil
}

func Defaults() []Func {
	return []Func{
		&Print{},
	}
}
