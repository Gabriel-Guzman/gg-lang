package operators

import "gg-lang/src/variables"

// string + string
type plusStrings struct{}

func (p *plusStrings) Evaluate(left, right interface{}) interface{} {
	return left.(string) + right.(string)
}

func (p *plusStrings) ResultType() variables.VarType {
	return variables.String
}
