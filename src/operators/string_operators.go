package operators

import (
	"gg-lang/src/variables"
	"strconv"
)

// string + string
type plusStrings struct{}

func (p *plusStrings) Evaluate(left, right interface{}) interface{} {
	return left.(string) + right.(string)
}

func (p *plusStrings) ResultType() variables.VarType {
	return variables.String
}

type intPlusString struct{}

func (*intPlusString) Evaluate(left, right interface{}) interface{} {
	return strconv.Itoa(left.(int)) + right.(string)
}

func (*intPlusString) ResultType() variables.VarType {
	return variables.String
}

type stringPlusInt struct{}

func (*stringPlusInt) Evaluate(left, right interface{}) interface{} {
	return left.(string) + strconv.Itoa(right.(int))
}

func (*stringPlusInt) ResultType() variables.VarType {
	return variables.String
}
