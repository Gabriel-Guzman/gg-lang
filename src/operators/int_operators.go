package operators

import (
	"gg-lang/src/variable"
)

// int + int
type plusInts struct{}

func (p *plusInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) + right.(int)
}

func (p *plusInts) ResultType() variable.VarType {
	return variable.Integer
}

// int - int
type minusInts struct{}

func (m *minusInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) - right.(int)
}

func (m *minusInts) ResultType() variable.VarType {
	return variable.Integer
}

func (m *minusInts) UnaryEvaluate(right interface{}) interface{} {
	return -right.(int)
}

// int * int
type mulInts struct{}

func (m *mulInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) * right.(int)
}

func (m *mulInts) ResultType() variable.VarType {
	return variable.Integer
}

// int / int
type divInts struct{}

func (d *divInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) / right.(int)
}

func (d *divInts) ResultType() variable.VarType {
	return variable.Integer
}

// int < int
type lessThanInts struct{}

func (l *lessThanInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) < right.(int)
}
func (l *lessThanInts) ResultType() variable.VarType { return variable.Boolean }

// int > int
type greaterThanInts struct{}

func (g *greaterThanInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) > right.(int)
}
func (g *greaterThanInts) ResultType() variable.VarType { return variable.Boolean }

// int <= int
type lessThanEqualInts struct{}

func (l *lessThanEqualInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) <= right.(int)
}
func (l *lessThanEqualInts) ResultType() variable.VarType { return variable.Boolean }

// int >= int
type greaterThanEqualInts struct{}

func (g *greaterThanEqualInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) >= right.(int)
}
func (g *greaterThanEqualInts) ResultType() variable.VarType { return variable.Boolean }

// -int
type minusInt struct{}

func (m *minusInt) Evaluate(right interface{}) interface{} {
	return -right.(int)
}
func (m *minusInt) ResultType() variable.VarType { return variable.Integer }
