package main

// string + string
type plusStrings struct{}

func (p *plusStrings) evaluate(left, right interface{}) interface{} {
	return left.(string) + right.(string)
}

func (p *plusStrings) resultType() varType {
	return STRING
}
