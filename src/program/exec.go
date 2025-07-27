package program

import (
	"gg-lang/src/gg"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variable"
)

func (p *Program) RunExpression(expr gg_ast.Expression) error {
	// dont execute anything if there's a return value right now
	if p.returnValue != nil {
		return nil
	}
	switch expr.(type) {
	case *gg_ast.TryCatchExpression:
		expr := expr.(*gg_ast.TryCatchExpression)
		if err := p.evaluateTryCatchExpression(expr); err != nil {
			return err
		}
	case *gg_ast.ReturnStatement:
		val, err := p.evaluateValueExpr(expr.(*gg_ast.ReturnStatement).Value)
		if err != nil {
			return err
		}
		p.returnValue = val
	case gg_ast.BlockStatement:
		block := expr.(gg_ast.BlockStatement)
		err := p.runBlockStmtNewScope(block)
		if err != nil {
			return err
		}
	case *gg_ast.AssignmentExpression:
		if err := p.evaluateAssignment(expr.(*gg_ast.AssignmentExpression)); err != nil {
			return err
		}
	case *gg_ast.DotAccessAssignmentExpression:
		if err := p.evaluateDotAccessAssignment(expr.(*gg_ast.DotAccessAssignmentExpression)); err != nil {
			return err
		}
	case *gg_ast.FunctionDeclExpression:
		decl := expr.(*gg_ast.FunctionDeclExpression)
		_, err := p.currentScope().declareVar(decl.Target.Tok.Symbol, &variable.RuntimeValue{
			Val: NewRuntimeFunc(decl, p.currentScope()),
			Typ: variable.Function,
		})
		if err != nil {
			return err
		}
	case *gg_ast.ForLoopExpression:
		err, done := p.execForLoopExpression(expr.(*gg_ast.ForLoopExpression))
		if done {
			return err
		}
	case *gg_ast.IfElseStatement:
		ifElse := expr.(*gg_ast.IfElseStatement)
		err := p.execIfElse(ifElse)
		if err != nil {
			return err
		}
	case *gg_ast.FunctionCallExpression:
		call := expr.(*gg_ast.FunctionCallExpression)
		_, err := p.call(call)
		if err != nil {
			return err
		}
	default:
		return gg.Crit("Invalid top-level expression: %s\n%s", expr.Kind().String(), gg_ast.NoBuilderExprString(expr))
	}
	return nil
}

func (p *Program) execForLoopExpression(expr *gg_ast.ForLoopExpression) (error, bool) {
	loop := expr
	for {
		val, err := p.evaluateValueExpr(loop.Condition)
		if err != nil {
			return err, true
		}
		if _, ok := val.Val.(bool); !ok {
			return gg.Runtime("loop condition must evaluate to bool\n%+v", expr), true
		}
		if !val.Val.(bool) {
			break
		}
		err = p.runBlockStmtNewScope(loop.Body)
		if err != nil {
			return err, true
		}
	}
	return nil, false
}

func (p *Program) execIfElse(expr gg_ast.Expression) error {
	ifElse := expr.(*gg_ast.IfElseStatement)
	cond, err := p.evaluateValueExpr(ifElse.Condition)
	if err != nil {
		return err
	}
	if _, ok := cond.Val.(bool); !ok {
		return gg.Runtime("if condition must evaluate to bool\n%+v", expr)
	}
	if cond.Val.(bool) {
		err = p.runBlockStmtNewScope(ifElse.Body)
		if err != nil {
			return err
		}
	} else {
		if ifElse.ElseExpression != nil {
			err = p.RunExpression(ifElse.ElseExpression)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Program) evaluateAssignment(expr *gg_ast.AssignmentExpression) error {
	if expr.Target.Kind() != gg_ast.ExprVariable {
		return gg.Runtime("invalid assignment target: %s", expr.Target.Tok.Symbol)
	}

	if expr.Value.Kind() > gg_ast.SentinelValueExpression {
		return gg.Runtime("cannot make value for %v", expr)
	}

	val, err := p.evaluateValueExpr(expr.Value)
	if err != nil {
		return err
	}

	existing := p.findVariable(expr.Target.Name())
	if existing != nil {
		existing.RuntimeValue = val // garbage collect old value
		return nil
	}

	_, err = p.currentScope().softDeclareVar(expr.Target.Tok.Symbol, val)
	return err
}
