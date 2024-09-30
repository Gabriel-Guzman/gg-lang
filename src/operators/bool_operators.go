package operators

import "gg-lang/src/variables"

type equalsAlwaysTrue struct{}

func (e *equalsAlwaysTrue) Evaluate(left, right interface{}) interface{} {
	return true
}
func (e *equalsAlwaysTrue) ResultType() variables.VarType {
	return variables.Boolean
}

type equalsAlwaysFalse struct{}

func (e *equalsAlwaysFalse) Evaluate(left, right interface{}) interface{} {
	return false
}
func (e *equalsAlwaysFalse) ResultType() variables.VarType {
	return variables.Boolean
}

type genEquals struct{}

func (g *genEquals) Evaluate(left, right interface{}) interface{} {
	return left == right
}

func (g *genEquals) ResultType() variables.VarType {
	return variables.Boolean
}

type genNotEquals struct{}

func (g *genNotEquals) Evaluate(left, right interface{}) interface{} {
	return left != right
}
func (g *genNotEquals) ResultType() variables.VarType {
	return variables.Boolean
}

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

type equalsBools struct{}

func (e *equalsBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) == rhs.(bool)
}
func (e *equalsBools) ResultType() variables.VarType {
	return variables.Boolean
}

type notEqualsBools struct{}

func (n *notEqualsBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) != rhs.(bool)
}
func (n *notEqualsBools) ResultType() variables.VarType {
	return variables.Boolean
}
