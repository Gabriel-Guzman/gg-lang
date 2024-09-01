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
