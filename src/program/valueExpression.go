package program

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/godTree"
	"gg-lang/src/variables"
	"strconv"
)

func (p *Program) evaluateValueExpr(expr godTree.IValExpr) (*variables.RuntimeValue, error) {
	switch expr.Kind() {
	case godTree.ExprVariable:
		name := expr.(*godTree.Identifier).Name()
		v := p.findVariable(name)
		if v != nil {
			return v.Value, nil
		}
		return nil, ggErrs.Runtime("undefined variable: %s", name)
	case godTree.ExprNumberLiteral:
		name := expr.(*godTree.Identifier).Name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return nil, ggErrs.Crit("unable to evaluate number literal: %s", err.Error())
		}
		return &variables.RuntimeValue{
			Val: intVal,
			Typ: variables.Integer,
		}, nil
	case godTree.ExprStringLiteral:
		return &variables.RuntimeValue{
			Val: expr.(*godTree.Identifier).Name(),
			Typ: variables.String,
		}, nil
	case godTree.ExprBinary:
		binExp := expr.(*godTree.BinaryExpression)

		left, err := p.evaluateValueExpr(binExp.Lhs)
		if err != nil {
			return nil, err
		}

		right, err := p.evaluateValueExpr(binExp.Rhs)
		if err != nil {
			return nil, err
		}

		op, exists := p.opMap.Get(binExp.Op, left.Typ, right.Typ)
		if !exists {
			return nil, ggErrs.Runtime("evaluateValueExpr: op %s not supported between types %T and %T", binExp.Op, left.Typ, right.Typ)
		}

		value := op.Evaluate(left.Val, right.Val)
		resultType := op.ResultType()

		return &variables.RuntimeValue{
			Val: value,
			Typ: resultType,
		}, nil
	case godTree.ExprFunctionCall:
		f := expr.(*godTree.FunctionCallExpression)
		if err2 := p.funcCall(f); err2 != nil {
			return &variables.RuntimeValue{
				Val: nil,
				Typ: variables.Void,
			}, nil
		} else {
			return nil, err2
		}
	default:
		return nil, ggErrs.Crit("evaluateValueExpr: invalid expression type: %v", expr)
	}
}
