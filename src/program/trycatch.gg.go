package program

import (
	"gg-lang/src/gg"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variable"
)

func (p *Program) evaluateCatchExpression(expr *gg_ast.CatchExpression, thrownErr *gg.RuntimeErr) error {
	p.enterNewScope()
	defer p.exitScope()

	_, err := p.currentScope().declareVar(expr.ErrorParam, &variable.RuntimeValue{Val: thrownErr.Error(), Typ: variable.String})
	if err != nil {
		return err
	}
	p.enterNewScope()
	defer p.exitScope()

	err = p.runBlockStmt(*expr.Body)
	if err != nil {
		return err
	}
	return nil
}

func (p *Program) evaluateTryCatchExpression(expr *gg_ast.TryCatchExpression) error {
	tryBlock := expr.Try
	catchBlock := expr.Catch
	finallyBlock := expr.Finally

	err := p.runBlockStmtNewScope(*tryBlock)
	if err != nil {
		if _, ok := err.(*gg.RuntimeErr); ok {
			err = p.evaluateCatchExpression(catchBlock, err.(*gg.RuntimeErr))
			if err != nil {
				return err
			}

			if finallyBlock != nil {
				err = p.runBlockStmtNewScope(*finallyBlock)
				if err != nil {
					return err
				}
			}

			return nil
		}
		return err
	}

	if finallyBlock != nil {
		err = p.runBlockStmtNewScope(*finallyBlock)
		if err != nil {
			return err
		}
	}

	return nil
}
