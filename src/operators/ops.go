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

	opm.Set("+", variables.Integer, variables.Integer, &plusInts{})
	opm.Set("-", variables.Integer, variables.Integer, &minusInts{})
	opm.Set("*", variables.Integer, variables.Integer, &mulInts{})
	opm.Set("/", variables.Integer, variables.Integer, &divInts{})

	opm.Set("+", variables.String, variables.String, &plusStrings{})
	opm.Set("+", variables.Integer, variables.String, &intPlusString{})
	opm.Set("+", variables.String, variables.Integer, &stringPlusInt{})

	opm.Set("==", variables.Boolean, variables.Boolean, &equalsBools{})
	opm.Set("!=", variables.Boolean, variables.Boolean, &notEqualsBools{})
	opm.Set("&&", variables.Boolean, variables.Boolean, &andBools{})
	opm.Set("||", variables.Boolean, variables.Boolean, &orBools{})

	opm.Set("==", variables.String, variables.String, &genEquals{})
	opm.Set("==", variables.Integer, variables.Integer, &genEquals{})

	return opm
}

var PrecedenceMap = map[string]int{
	"*":  2,
	"/":  2,
	"+":  1,
	"-":  1,
	"&&": 0,
	"||": 0,
	"==": -1,
	"!=": -1,
}

func LeftFirst(l, r string) bool {
	pl, ok := PrecedenceMap[l]
	if !ok {
		panic(fmt.Sprintf("checked precedence on nonexistent op %s", l))
	}
	pr, ok := PrecedenceMap[r]
	if !ok {
		panic(fmt.Sprintf("checked precedence on nonexistent op %s", r))
	}

	if pl == pr {
		return true
	}

	return pl > pr
}
