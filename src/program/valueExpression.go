package program

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/godTree"
	"gg-lang/src/variables"
	"strconv"
)

func (p *Program) evaluateValueExpr(expr godTree.ValueExpression) (interface{}, variables.VarType, error) {
	switch expr.Kind() {
	case godTree.ExprVariable:
		name := expr.(*godTree.Identifier).Name()
		v := p.findVariable(name)
		if v != nil {
			return v.Value, v.Typ, nil
		}
		return nil, 0, ggErrs.Runtime("undefined variable: %s", name)
	case godTree.ExprNumberLiteral:
		name := expr.(*godTree.Identifier).Name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return nil, 0, ggErrs.Crit("unable to evaluate number literal: %s", err.Error())
		}
		return intVal, variables.Integer, nil
	case godTree.ExprStringLiteral:
		return expr.(*godTree.Identifier).Name(), variables.String, nil
	case godTree.ExprBinary:
		binExp := expr.(*godTree.BinaryExpression)
		left, ltyp, err := p.evaluateValueExpr(binExp.Lhs)
		if err != nil {
			return nil, 0, err
		}

		right, rtyp, err := p.evaluateValueExpr(binExp.Rhs)
		if err != nil {
			return nil, 0, err
		}

		op, exists := p.opMap.Get(binExp.Op, ltyp, rtyp)
		if !exists {
			return nil, 0, ggErrs.Runtime("evaluateValueExpr: op %s not supported between types %v and %v", binExp.Op, ltyp, rtyp)
		}

		value := op.Evaluate(left, right)
		resultType := op.ResultType()

		return value, resultType, nil
	case godTree.ExprFunctionCall:
		f := expr.(*godTree.FunctionCallExpression)
		if err := p.funcCall(f); err != nil {
			return 0, variables.Integer, err
		}

		return 0, variables.Integer, nil
	default:
		return nil, 0, ggErrs.Crit("evaluateValueExpr: invalid expression type: %v", expr)
	}
}
