package gg_ast

import (
	"gg-lang/src/variable"
	"strconv"
)

// string + string
type plusStrings struct{}

func (p *plusStrings) Evaluate(left, right interface{}) interface{} {
	return left.(string) + right.(string)
}

func (p *plusStrings) ResultType() variable.VarType {
	return variable.String
}

type intPlusString struct{}

func (*intPlusString) Evaluate(left, right interface{}) interface{} {
	return strconv.Itoa(left.(int)) + right.(string)
}

func (*intPlusString) ResultType() variable.VarType {
	return variable.String
}

type stringPlusInt struct{}

func (*stringPlusInt) Evaluate(left, right interface{}) interface{} {
	return left.(string) + strconv.Itoa(right.(int))
}

func (*stringPlusInt) ResultType() variable.VarType {
	return variable.String
}

type coercedPlusString struct{}

func (*coercedPlusString) Evaluate(left, right interface{}) interface{} {
	lhs, err := variable.CoerceTo(left, variable.String)
	if err != nil {
		return nil
	}
	return lhs.(string) + right.(string)
}

func (*coercedPlusString) ResultType() variable.VarType {
	return variable.String
}

type stringPlusCoerced struct{}

func (*stringPlusCoerced) Evaluate(left, right interface{}) interface{} {
	rhs, err := variable.CoerceTo(right, variable.String)
	if err != nil {
		return nil
	}
	return left.(string) + rhs.(string)
}
func (*stringPlusCoerced) ResultType() variable.VarType {
	return variable.String
}
