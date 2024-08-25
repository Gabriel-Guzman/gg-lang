package operators

import (
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src/variables"
)

type Operator interface {
	Evaluate(left, right interface{}) interface{}
	ResultType() variables.VarType
}

type Opmap struct {
	ops map[string]Operator
}

func opKey(name string, left, right variables.VarType) string {
	return fmt.Sprintf("%s_%d_%d", name, left, right)
}

func (o *Opmap) Get(name string, left, right variables.VarType) (Operator, bool) {
	op, ok := o.ops[opKey(name, left, right)]
	return op, ok
}

func (o *Opmap) set(name string, left, right variables.VarType, op Operator) {
	o.ops[opKey(name, left, right)] = op
}

func DefaultOpMap() *Opmap {
	opm := &Opmap{
		ops: make(map[string]Operator),
	}

	plus := plusInts{}
	opm.set("+", variables.INTEGER, variables.INTEGER, &plus)

	minus := minusInts{}
	opm.set("-", variables.INTEGER, variables.INTEGER, &minus)

	mul := mulInts{}
	opm.set("*", variables.INTEGER, variables.INTEGER, &mul)

	div := divInts{}
	opm.set("/", variables.INTEGER, variables.INTEGER, &div)

	plusStr := plusStrings{}
	opm.set("+", variables.STRING, variables.STRING, &plusStr)

	return opm
}
