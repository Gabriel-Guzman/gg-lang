package operators

import "gg-lang/src/variables"

type andBools struct{}

func (a *andBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) && rhs.(bool)
}
func (a *andBools) ResultType() variables.VarType {
	return variables.Boolean
}

type orBools struct{}

func (o *orBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) || rhs.(bool)
}
func (o *orBools) ResultType() variables.VarType {
	return variables.Boolean
}
