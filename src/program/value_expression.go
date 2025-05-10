package program

import (
	"gg-lang/src/gg"
	"gg-lang/src/gg_ast"
	"gg-lang/src/operators"
	"gg-lang/src/variable"
	"strconv"
)

func (p *Program) evaluateValueExpr(expr gg_ast.ValueExpression) (*variable.RuntimeValue, error) {
	switch expr.Kind() {
	case gg_ast.ExprVariable:
		name := expr.(*gg_ast.Identifier).Name()
		v := p.findVariable(name)
		if v != nil {
			return v.RuntimeValue, nil
		}
		return nil, gg.Runtime("undefined variable: %s", name)
	case gg_ast.ExprIntLiteral:
		name := expr.(*gg_ast.Identifier).Name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return nil, gg.Crit("unable to evaluate int literal: %s", err.Error())
		}
		return &variable.RuntimeValue{
			Val: intVal,
			Typ: variable.Integer,
		}, nil
	case gg_ast.ExprBoolLiteral:
		name := expr.(*gg_ast.Identifier).Name()
		boolVal, err := strconv.ParseBool(name)
		if err != nil {
			return nil, gg.Crit("unable to evaluate bool literal: %s", err.Error())
		}
		return &variable.RuntimeValue{
			Val: boolVal,
			Typ: variable.Boolean,
		}, nil
	case gg_ast.ExprStringLiteral:
		return &variable.RuntimeValue{
			Val: expr.(*gg_ast.Identifier).Name(),
			Typ: variable.String,
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

		op, exists := operators.Get(binExp.Op.Symbol, left.Typ, right.Typ)
		if !exists {
			return nil, gg.Runtime(
				"evaluateValueExpr: op %s not supported between types %s and %s\nevaluating: %s", binExp.Op, left.Typ.String(), right.Typ.String(), gg_ast.NoBuilderExprString(expr))
		}

		value := op.Evaluate(left.Val, right.Val)
		resultType := op.ResultType()

		return &variable.RuntimeValue{
			Val: value,
			Typ: resultType,
		}, nil
	case gg_ast.ExprFunctionCall:
		f := expr.(*gg_ast.FunctionCallExpression)
		return p.call(f)
	case gg_ast.ExprFuncDecl:
		decl := expr.(*gg_ast.FunctionDeclExpression)
		return &variable.RuntimeValue{
			Val: NewRuntimeFunc(decl, p.currentScope()),
			Typ: variable.Function,
		}, nil
	case gg_ast.ExprUnary:
		e := expr.(*gg_ast.UnaryExpression)
		rhs, err := p.evaluateValueExpr(e.Rhs)
		if err != nil {
			return nil, err
		}
		op, exists := operators.GetUnary(e.Op.Symbol, rhs.Typ)
		if !exists {
			return nil, gg.Runtime(
				"evaluateValueExpr: unary op %s not supported for type %s\nevaluating: %s", e.Op.Symbol, rhs.Typ.String(), gg_ast.NoBuilderExprString(expr))
		}
		value := op.Evaluate(rhs.Val)
		return &variable.RuntimeValue{
			Val: value,
			Typ: rhs.Typ,
		}, nil

	default:
		return nil, gg.Crit("evaluateValueExpr: invalid expression type: %v", expr)
	}
}
