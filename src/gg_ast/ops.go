package gg_ast

import (
	"fmt"
	"gg-lang/src/variable"
	"strings"
)

type Operator interface {
	Evaluate(left, right interface{}) interface{}
	ResultType() variable.VarType
}

type OpMap struct {
	ops map[string]Operator
}

func opKey(name string, left, right variable.VarType) string {
	return fmt.Sprintf("%s_%d_%d", name, left, right)
}

func (o *OpMap) Get(name string, left, right variable.VarType) (Operator, bool) {
	op, ok := o.ops[opKey(name, left, right)]
	return op, ok
}

func (o *OpMap) Set(name string, left, right variable.VarType, op Operator) {
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

	opm.Set("+", variable.Integer, variable.Integer, &plusInts{})
	opm.Set("-", variable.Integer, variable.Integer, &minusInts{})
	opm.Set("*", variable.Integer, variable.Integer, &mulInts{})
	opm.Set("/", variable.Integer, variable.Integer, &divInts{})

	opm.Set("<", variable.Integer, variable.Integer, &lessThanInts{})
	opm.Set(">", variable.Integer, variable.Integer, &greaterThanInts{})
	opm.Set("<=", variable.Integer, variable.Integer, &lessThanEqualInts{})
	opm.Set(">=", variable.Integer, variable.Integer, &greaterThanEqualInts{})
	opm.Set("!=", variable.Integer, variable.Integer, &genNotEquals{})
	opm.Set("==", variable.Integer, variable.Integer, &genEquals{})

	opm.Set("+", variable.String, variable.String, &plusStrings{})
	opm.Set("+", variable.Integer, variable.String, &coercedPlusString{})
	opm.Set("+", variable.String, variable.Integer, &stringPlusCoerced{})
	opm.Set("+", variable.Boolean, variable.String, &coercedPlusString{})
	opm.Set("+", variable.String, variable.Boolean, &stringPlusCoerced{})

	opm.Set("==", variable.Boolean, variable.Boolean, &equalsBools{})
	opm.Set("!=", variable.Boolean, variable.Boolean, &notEqualsBools{})
	opm.Set("&&", variable.Boolean, variable.Boolean, &andBools{})
	opm.Set("||", variable.Boolean, variable.Boolean, &orBools{})

	opm.Set("==", variable.String, variable.String, &genEquals{})
	opm.Set("!=", variable.String, variable.String, &genNotEquals{})

	opm.Set("==", variable.Void, variable.Void, &equalsAlwaysTrue{})
	opm.Set("!=", variable.Void, variable.Void, &equalsAlwaysFalse{})
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
	"<":  -1,
	">":  -1,
	"<=": -1,
	">=": -1,
}
