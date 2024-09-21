package operators

import "gg-lang/src/variables"

// int + int
type plusInts struct{}

func (p *plusInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) + right.(int)
}

func (p *plusInts) ResultType() variables.VarType {
	return variables.Integer
}

// int - int
type minusInts struct{}

func (m *minusInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) - right.(int)
}

func (m *minusInts) ResultType() variables.VarType {
	return variables.Integer
}

// int * int
type mulInts struct{}

func (m *mulInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) * right.(int)
}

func (m *mulInts) ResultType() variables.VarType {
	return variables.Integer
}

// int / int
type divInts struct{}

func (d *divInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) / right.(int)
}

func (d *divInts) ResultType() variables.VarType {
	return variables.Integer
}

// int < int
type lessThanInts struct{}

func (l *lessThanInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) < right.(int)
}
func (l *lessThanInts) ResultType() variables.VarType { return variables.Boolean }

// int > int
type greaterThanInts struct{}

func (g *greaterThanInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) > right.(int)
}
func (g *greaterThanInts) ResultType() variables.VarType { return variables.Boolean }

// int <= int
type lessThanEqualInts struct{}

func (l *lessThanEqualInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) <= right.(int)
}
func (l *lessThanEqualInts) ResultType() variables.VarType { return variables.Boolean }

// int >= int
type greaterThanEqualInts struct{}

func (g *greaterThanEqualInts) Evaluate(left, right interface{}) interface{} {
	return left.(int) >= right.(int)
}
func (g *greaterThanEqualInts) ResultType() variables.VarType { return variables.Boolean }
