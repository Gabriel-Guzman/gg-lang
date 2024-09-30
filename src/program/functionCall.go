package program

import (
	"gg-lang/src/builtin"
	"gg-lang/src/ggErrs"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variables"
)

func (p *Program) call(f *gg_ast.FunctionCallExpression) (*variables.RuntimeValue, error) {
	variable, ok := p.top.variables[f.Id.Raw]
	if !ok {
		return nil, ggErrs.Runtime("undefined function %s, evaluating\n%s", f.Id.Raw, gg_ast.NoBuilderExprString(f))
	}

	// check if callable
	if variable.RuntimeValue.Typ != variables.Function &&
		variable.RuntimeValue.Typ != variables.BuiltinFunction {
		return nil, ggErrs.Runtime("%s is not callable, evaluating\n%s", f.Id.Raw, gg_ast.NoBuilderExprString(f))
	}

	// build values for arguments
	vals := make([]*variables.RuntimeValue, len(f.Args))
	for i, arg := range f.Args {
		value, err := p.evaluateValueExpr(arg)
		if err != nil {
			return nil, err
		}

		vals[i] = value
	}

	// run builtin
	if variable.RuntimeValue.Typ == variables.BuiltinFunction {
		return p.builtinFuncCall(variable.RuntimeValue.Val.(builtin.Func), vals)
	}

	// set up func expression
	funcDeclExpr := variable.RuntimeValue.Val.(*gg_ast.FunctionDeclExpression)
	if len(funcDeclExpr.Params) != len(f.Args) {
		return nil, ggErrs.Runtime("param count mismatch on %s, evaluating\n%s", f.Id.Raw, gg_ast.NoBuilderExprString(f))
	}

	// build variables for new scope
	var scopedVariables []variables.Variable
	for i := range f.Args {
		value := vals[i]

		name := funcDeclExpr.Params[i]
		scopedVariables = append(scopedVariables, variables.Variable{
			Name:         name,
			RuntimeValue: value,
		})
	}

	p.enterNewScope()
	defer p.exitScope()

	for _, v := range scopedVariables {
		temp := v
		p.current.variables[v.Name] = &temp
	}

	for _, stmt := range funcDeclExpr.Body {
		if stmt.Kind() == gg_ast.ExprReturn {
			return p.evaluateValueExpr(stmt.(*gg_ast.ReturnStatement).Value)
		}
		err := p.RunStmt(stmt)
		if err != nil {
			return nil, err
		}
	}

	return &variables.RuntimeValue{
		Val: nil,
		Typ: variables.Void,
	}, nil
}

func (p *Program) builtinFuncCall(f builtin.Func, args []*variables.RuntimeValue) (*variables.RuntimeValue, error) {
	//var vals []*variables.RuntimeValue
	//for _, arg := range args {
	//	vals = append(vals, arg.Val)
	//}
	res, err := f.Call(args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
