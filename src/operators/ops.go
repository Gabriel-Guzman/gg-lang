package operators

import (
	"fmt"
	"gg-lang/src/variables"
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
		sb.WriteString(fmt.Sprintf("\t%s: %T\n", key, op))
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

	opm.Set("+", variables.Integer, variables.String, &intPlusString{})
	opm.Set("+", variables.String, variables.Integer, &stringPlusInt{})

	return opm
}

var PrecedenceMap map[string]int = map[string]int{
	"*": 2,
	"/": 2,
	"+": 1,
	"-": 1,
}

func LeftFirst(l, r string) bool {
	pl, ok := PrecedenceMap[l]
	if !ok {
		panic(fmt.Sprintf("checked precedence on nonexistent op %s or %s", l, r))
	}
	rl, ok := PrecedenceMap[r]
	if !ok {
		panic(fmt.Sprintf("checked precedence on nonexistent op %s or %s", l, r))
	}

	if pl == rl {
		return true
	}

	return pl > rl
}

//func PrecedenceMap() map[string]int {
//	return map[string]int{
//		"+": 1,
//		"-": 1,
//		"*": 2,
//		"/": 2,
//	}
//}
