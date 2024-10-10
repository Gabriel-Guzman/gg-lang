package program

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variable"
)

func (p *Program) call(f *gg_ast.FunctionCallExpression) (*variable.RuntimeValue, error) {
	v := p.currentScope().findVariable(f.Name())
	if v == nil {
		return nil, ggErrs.Runtime("undefined function %s, evaluating\n%s", f.Id.Raw, gg_ast.NoBuilderExprString(f))
	}

	// check if callable
	if v.RuntimeValue.Typ != variable.Function &&
		v.RuntimeValue.Typ != variable.BuiltinFunction {
		return nil, ggErrs.Runtime("%s is not callable, evaluating\n%s", f.Id.Raw, gg_ast.NoBuilderExprString(f))
	}

	// build values for arguments
	vals := make([]*variable.RuntimeValue, len(f.Args))
	for i, arg := range f.Args {
		value, err := p.evaluateValueExpr(arg)
		if err != nil {
			return nil, err
		}

		vals[i] = value
	}

	// run builtin
	if bn, ok := v.RuntimeValue.Val.(Func); ok {
		return p.builtinFuncCall(bn, vals)
	}

	// set up func expression
	runtimeFunc := v.RuntimeValue.Val.(*RuntimeFunc)
	if len(runtimeFunc.Decl.Params) != len(f.Args) {
		return nil, ggErrs.Runtime("param count mismatch on %s, evaluating\n%s", f.Id.Raw, gg_ast.NoBuilderExprString(f))
	}

	// build variables for new scope
	var scopedVariables []variable.Variable
	for i := range f.Args {
		value := vals[i]

		name := runtimeFunc.Decl.Params[i]
		scopedVariables = append(scopedVariables, variable.Variable{
			Name:         name,
			RuntimeValue: value,
		})
	}
	// enter the captured scope, and then a new one to temporarily save the arguments
	p.enterCapturedScope(runtimeFunc.CapturedScope)
	p.enterNewScope()
	defer p.exitScope()
	defer p.exitScope()

	for _, v := range scopedVariables {
		temp := v
		p.currentScope().variables[temp.Name] = &temp
	}
	for _, stmt := range runtimeFunc.Decl.Body {
		err := p.RunStmt(stmt)
		if err != nil {
			return nil, err
		}
		if p.returnValue != nil {
			ret := p.returnValue
			p.returnValue = nil
			return ret, nil
		}
	}
	return &variable.RuntimeValue{
		Typ: variable.Void,
	}, nil
}

func (p *Program) builtinFuncCall(f Func, args []*variable.RuntimeValue) (*variable.RuntimeValue, error) {
	res, err := f.Call(args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
