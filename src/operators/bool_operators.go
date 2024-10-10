package operators

import (
	"gg-lang/src/variable"
)

type equalsAlwaysTrue struct{}

func (e *equalsAlwaysTrue) Evaluate(left, right interface{}) interface{} {
	return true
}
func (e *equalsAlwaysTrue) ResultType() variable.VarType {
	return variable.Boolean
}

type equalsAlwaysFalse struct{}

func (e *equalsAlwaysFalse) Evaluate(left, right interface{}) interface{} {
	return false
}
func (e *equalsAlwaysFalse) ResultType() variable.VarType {
	return variable.Boolean
}

type genEquals struct{}

func (g *genEquals) Evaluate(left, right interface{}) interface{} {
	return left == right
}

func (g *genEquals) ResultType() variable.VarType {
	return variable.Boolean
}

type genNotEquals struct{}

func (g *genNotEquals) Evaluate(left, right interface{}) interface{} {
	return left != right
}
func (g *genNotEquals) ResultType() variable.VarType {
	return variable.Boolean
}

type andBools struct{}

func (a *andBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) && rhs.(bool)
}
func (a *andBools) ResultType() variable.VarType {
	return variable.Boolean
}

type orBools struct{}

func (o *orBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) || rhs.(bool)
}
func (o *orBools) ResultType() variable.VarType {
	return variable.Boolean
}

type equalsBools struct{}

func (e *equalsBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) == rhs.(bool)
}
func (e *equalsBools) ResultType() variable.VarType {
	return variable.Boolean
}

type notEqualsBools struct{}

func (n *notEqualsBools) Evaluate(lhs interface{}, rhs interface{}) interface{} {
	return lhs.(bool) != rhs.(bool)
}
func (n *notEqualsBools) ResultType() variable.VarType {
	return variable.Boolean
}
