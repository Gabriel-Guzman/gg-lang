package program

import "gg-lang/src/gg_ast"

type RuntimeFunc struct {
	Name          string
	Decl          *gg_ast.FunctionDeclExpression
	CapturedScope *Scope
}

func RuntimeFuncFromDecl(decl *gg_ast.FunctionDeclExpression, scope *Scope) *RuntimeFunc {
	return &RuntimeFunc{
		Name:          decl.Target.Name(),
		Decl:          decl,
		CapturedScope: scope,
	}
}
