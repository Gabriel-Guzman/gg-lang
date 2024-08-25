package main

import "fmt"

type operator interface {
	evaluate(left, right interface{}) interface{}
	resultType() varType
}

type opmap struct {
	ops map[string]operator
}

func opKey(name string, left, right varType) string {
	return fmt.Sprintf("%s_%d_%d", name, left, right)
}

func (o *opmap) get(name string, left, right varType) (operator, bool) {
	op, ok := o.ops[opKey(name, left, right)]
	return op, ok
}

func (o *opmap) set(name string, left, right varType, op operator) {
	o.ops[opKey(name, left, right)] = op
}

func defaultOpMap() *opmap {
	opm := &opmap{
		ops: make(map[string]operator),
	}

	plus := plusInts{}
	opm.set("+", INTEGER, INTEGER, &plus)

	minus := minusInts{}
	opm.set("-", INTEGER, INTEGER, &minus)

	mul := mulInts{}
	opm.set("*", INTEGER, INTEGER, &mul)

	div := divInts{}
	opm.set("/", INTEGER, INTEGER, &div)

	plusStr := plusStrings{}
	opm.set("+", STRING, STRING, &plusStr)

	return opm
}
