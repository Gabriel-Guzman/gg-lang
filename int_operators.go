package main

// int + int
type plusInts struct{}

func (p *plusInts) evaluate(left, right interface{}) interface{} {
	return left.(int) + right.(int)
}

func (p *plusInts) resultType() varType {
	return INTEGER
}

// int - int
type minusInts struct{}

func (m *minusInts) evaluate(left, right interface{}) interface{} {
	return left.(int) - right.(int)
}

func (m *minusInts) resultType() varType {
	return INTEGER
}

// int * int
type mulInts struct{}

func (m *mulInts) evaluate(left, right interface{}) interface{} {
	return left.(int) * right.(int)
}

func (m *mulInts) resultType() varType {
	return INTEGER
}

// int / int
type divInts struct{}

func (d *divInts) evaluate(left, right interface{}) interface{} {
	return left.(int) / right.(int)
}

func (d *divInts) resultType() varType {
	return INTEGER
}
