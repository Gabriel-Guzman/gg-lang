package program

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/godTree"
	"gg-lang/src/variables"
)

func (p *Program) funcCall(f *godTree.FunctionCallExpression) error {
	val, ok := p.top.variables[f.Id.Raw]
	if !ok {
		return ggErrs.Runtime("undefined function %s", f.Id.Raw)
	}
	if val.Typ != variables.Function {
		return ggErrs.Runtime("%s is not callable", f.Id.Raw)
	}

	funcDeclExpr := val.Value.(*godTree.FunctionDeclExpression)
	if len(funcDeclExpr.Parms) != len(f.Args) {
		return ggErrs.Runtime("param count mismatch on", f.Id.Raw)
	}

	var scopedVariables []variables.Variable
	for i, arg := range f.Args {
		value, typ, err := p.evaluateValueExpr(arg)
		if err != nil {
			return err
		}

		name := funcDeclExpr.Parms[i]

		scopedVariables = append(scopedVariables, variables.Variable{
			Name:  name,
			Typ:   typ,
			Value: value,
		})
	}
	p.enterNewScope()
	defer p.exitScope()

	for _, variable := range scopedVariables {
		temp := variable
		p.current.variables[variable.Name] = &temp
	}

	err := p.Run(&godTree.Ast{Body: funcDeclExpr.Value})
	if err != nil {
		return err
	}
	return nil
}
