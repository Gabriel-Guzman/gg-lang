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

type coercedPlusString struct{}

func (*coercedPlusString) Evaluate(left, right interface{}) interface{} {
	lhs, err := variables.CoerceTo(left, variables.String)
	if err != nil {
		return nil
	}
	return lhs.(string) + right.(string)
}

func (*coercedPlusString) ResultType() variables.VarType {
	return variables.String
}

type stringPlusCoerced struct{}

func (*stringPlusCoerced) Evaluate(left, right interface{}) interface{} {
	rhs, err := variables.CoerceTo(right, variables.String)
	if err != nil {
		return nil
	}
	return left.(string) + rhs.(string)
}
func (*stringPlusCoerced) ResultType() variables.VarType {
	return variables.String
}
