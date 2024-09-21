package program

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variables"
	"strconv"
)

func (p *Program) evaluateValueExpr(expr gg_ast.ValExpression) (*variables.RuntimeValue, error) {
	switch expr.Kind() {
	case gg_ast.ExprVariable:
		name := expr.(*gg_ast.Identifier).Name()
		v := p.findVariable(name)
		if v != nil {
			return v.RuntimeValue, nil
		}
		return nil, ggErrs.Runtime("undefined variable: %s", name)
	case gg_ast.ExprIntLiteral:
		name := expr.(*gg_ast.Identifier).Name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return nil, ggErrs.Crit("unable to evaluate int literal: %s", err.Error())
		}
		return &variables.RuntimeValue{
			Val: intVal,
			Typ: variables.Integer,
		}, nil
	case gg_ast.ExprBoolLiteral:
		name := expr.(*gg_ast.Identifier).Name()
		boolVal, err := strconv.ParseBool(name)
		if err != nil {
			return nil, ggErrs.Crit("unable to evaluate bool literal: %s", err.Error())
		}
		return &variables.RuntimeValue{
			Val: boolVal,
			Typ: variables.Boolean,
		}, nil
	case gg_ast.ExprStringLiteral:
		return &variables.RuntimeValue{
			Val: expr.(*gg_ast.Identifier).Name(),
			Typ: variables.String,
		}, nil
	case gg_ast.ExprBinary:
		binExp := expr.(*gg_ast.BinaryExpression)

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
	case gg_ast.ExprFunctionCall:
		f := expr.(*gg_ast.FunctionCallExpression)
		if err := p.call(f); err != nil {
			return &variables.RuntimeValue{
				Val: nil,
				Typ: variables.Void,
			}, nil
		} else {
			return nil, err
		}
	default:
		return nil, ggErrs.Crit("evaluateValueExpr: invalid expression type: %v", expr)
	}
}
