package operators

import (
	"fmt"
	"gg-lang/src/variable"
	"strings"
)

type Operator interface {
	Evaluate(left, right interface{}) interface{}
	ResultType() variable.VarType
}

type UnaryOperator interface {
	Evaluate(right interface{}) interface{}
	ResultType() variable.VarType
}

type UnaryOpMap struct {
	ops map[string]UnaryOperator
}

func (o *UnaryOpMap) setUnary(name string, right variable.VarType, op UnaryOperator) {
	o.ops[unaryOpKey(name, right)] = op
}

type OpMap struct {
	Ops map[string]Operator
}

func opKey(name string, left, right variable.VarType) string {
	return fmt.Sprintf("%s_%d_%d", name, left, right)
}

func unaryOpKey(name string, right variable.VarType) string {
	return fmt.Sprintf("%s_%d", name, right)
}

var defaultOpMap map[string]Operator
var defaultUnaryOpMap map[string]UnaryOperator

func init() {
	defaultOpMap = Default().Ops
	defaultUnaryOpMap = DefaultUnary().ops
}

func Get(name string, left, right variable.VarType) (Operator, bool) {
	ops := defaultOpMap
	op, ok := ops[opKey(name, left, right)]
	return op, ok
}

func GetUnary(name string, right variable.VarType) (UnaryOperator, bool) {
	op, ok := defaultUnaryOpMap[unaryOpKey(name, right)]
	uop, ok := op.(UnaryOperator)
	return uop, ok
}

func (o *OpMap) set(name string, left, right variable.VarType, op Operator) {
	o.Ops[opKey(name, left, right)] = op
}

func (o *OpMap) String() string {
	var sb strings.Builder
	for key, op := range o.Ops {
		sb.WriteString(fmt.Sprintf("\t%s: %T\n", key, op))
	}
	return sb.String()
}

func DefaultUnary() *UnaryOpMap {
	opm := &UnaryOpMap{
		ops: make(map[string]UnaryOperator),
	}

	opm.setUnary("-", variable.Integer, &minusInt{})
	opm.setUnary("!", variable.Boolean, &notBool{})

	return opm
}

func Default() *OpMap {
	opm := &OpMap{
		Ops: make(map[string]Operator),
	}

	opm.set("+", variable.Integer, variable.Integer, &plusInts{})
	opm.set("-", variable.Integer, variable.Integer, &minusInts{})
	opm.set("*", variable.Integer, variable.Integer, &mulInts{})
	opm.set("/", variable.Integer, variable.Integer, &divInts{})

	opm.set("<", variable.Integer, variable.Integer, &lessThanInts{})
	opm.set(">", variable.Integer, variable.Integer, &greaterThanInts{})
	opm.set("<=", variable.Integer, variable.Integer, &lessThanEqualInts{})
	opm.set(">=", variable.Integer, variable.Integer, &greaterThanEqualInts{})
	opm.set("!=", variable.Integer, variable.Integer, &genNotEquals{})
	opm.set("==", variable.Integer, variable.Integer, &genEquals{})

	opm.set("+", variable.String, variable.String, &plusStrings{})
	opm.set("+", variable.Integer, variable.String, &coercedPlusString{})
	opm.set("+", variable.String, variable.Integer, &stringPlusCoerced{})
	opm.set("+", variable.Boolean, variable.String, &coercedPlusString{})
	opm.set("+", variable.String, variable.Boolean, &stringPlusCoerced{})

	opm.set("==", variable.Boolean, variable.Boolean, &equalsBools{})
	opm.set("!=", variable.Boolean, variable.Boolean, &notEqualsBools{})
	opm.set("&&", variable.Boolean, variable.Boolean, &andBools{})
	opm.set("||", variable.Boolean, variable.Boolean, &orBools{})

	opm.set("==", variable.String, variable.String, &genEquals{})
	opm.set("!=", variable.String, variable.String, &genNotEquals{})

	opm.set("==", variable.Void, variable.Void, &equalsAlwaysTrue{})
	opm.set("!=", variable.Void, variable.Void, &equalsAlwaysFalse{})

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
