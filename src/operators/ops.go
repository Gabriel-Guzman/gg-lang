package operators

import (
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src/variables"
	"strings"
)

type Operator interface {
	Evaluate(left, right interface{}) interface{}
	ResultType() variables.VarType
}

type OpMap struct {
	ops map[string]Operator
}

func opKey(name string, left, right variables.VarType) string {
	return fmt.Sprintf("%s_%d_%d", name, left, right)
}

func (o *OpMap) Get(name string, left, right variables.VarType) (Operator, bool) {
	op, ok := o.ops[opKey(name, left, right)]
	return op, ok
}

func (o *OpMap) Set(name string, left, right variables.VarType, op Operator) {
	o.ops[opKey(name, left, right)] = op
}

func (o *OpMap) String() string {
	var sb strings.Builder
	for key, op := range o.ops {
		sb.WriteString(fmt.Sprintf("%s: %T\n", key, op))
	}
	return sb.String()
}

func Default() *OpMap {
	opm := &OpMap{
		ops: make(map[string]Operator),
	}

	plus := plusInts{}
	opm.Set("+", variables.Integer, variables.Integer, &plus)

	minus := minusInts{}
	opm.Set("-", variables.Integer, variables.Integer, &minus)

	mul := mulInts{}
	opm.Set("*", variables.Integer, variables.Integer, &mul)

	div := divInts{}
	opm.Set("/", variables.Integer, variables.Integer, &div)

	plusStr := plusStrings{}
	opm.Set("+", variables.String, variables.String, &plusStr)

	return opm
}
