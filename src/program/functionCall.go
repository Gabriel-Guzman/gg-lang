package program

import (
	"gg-lang/src/builtin"
	"gg-lang/src/ggErrs"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variables"
)

func (p *Program) funcCall(f *gg_ast.FunctionCallExpression) error {
	variable, ok := p.top.variables[f.Id.Raw]
	if !ok {
		return ggErrs.Runtime("undefined function %s", f.Id.Raw)
	}

	// check if callable
	if variable.Value.Typ != variables.Function &&
		variable.Value.Typ != variables.BuiltinFunction {
		return ggErrs.Runtime("%s is not callable", f.Id.Raw)
	}

	// build values for arguments
	vals := make([]*variables.RuntimeValue, len(f.Args))
	for i, arg := range f.Args {
		value, err := p.evaluateValueExpr(arg)
		if err != nil {
			return err
		}

		vals[i] = value
	}

	// run builtin
	if variable.Value.Typ == variables.BuiltinFunction {
		return p.builtinFuncCall(variable.Value.Val.(builtin.Func), vals)
	}

	// set up func expression
	funcDeclExpr := variable.Value.Val.(*gg_ast.FunctionDeclExpression)
	if len(funcDeclExpr.Params) != len(f.Args) {
		return ggErrs.Runtime("param count mismatch on", f.Id.Raw)
	}

	// build variables for new scope
	var scopedVariables []variables.Variable
	for i := range f.Args {
		value := vals[i]

		name := funcDeclExpr.Params[i]
		scopedVariables = append(scopedVariables, variables.Variable{
			Name:  name,
			Value: value,
		})
	}

	p.enterNewScope()
	defer p.exitScope()

	for _, v := range scopedVariables {
		temp := v
		p.current.variables[v.Name] = &temp
	}

	err := p.Run(&gg_ast.Ast{Body: funcDeclExpr.Value})
	if err != nil {
		return err
	}
	return nil
}

func (p *Program) builtinFuncCall(f builtin.Func, args []*variables.RuntimeValue) error {
	//var vals []*variables.RuntimeValue
	//for _, arg := range args {
	//	vals = append(vals, arg.Val)
	//}
	err := f.Call(args...)
	if err != nil {
		return ggErrs.Runtime("builtin function %s: %v", f.Name(), err)
	}

	return nil
}
