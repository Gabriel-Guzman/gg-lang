package program

import (
	"gg-lang/src/gg"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variable"
)

type Object = map[string]*variable.RuntimeValue

func (p *Program) getDotAccessAssignmentTarget(expr *gg_ast.DotAccessAssignmentExpression) (Object, error) {
	e := expr.Target
	var currentObject Object
	var res *variable.RuntimeValue

	for i, accessKey := range e.AccessChain[:len(e.AccessChain)-1] {
		if i == 0 {
			v := p.findVariable(accessKey)
			if v == nil {
				return nil, gg.Runtime("undefined variable: %s\nevaluating %s\n%s", accessKey, e.Name(), gg_ast.NoBuilderExprString(expr))
			}
			res = v.RuntimeValue
		} else {
			var ok bool
			res, ok = currentObject[accessKey]
			if !ok {
				return nil, gg.Runtime("undefined property: %s\nevaluating %s\n%s", accessKey, e.Name(), gg_ast.NoBuilderExprString(expr))
			}
		}

		if res.Typ != variable.Object {
			return nil, gg.Runtime("%s is not an object, evaluating\n%s", accessKey, gg_ast.NoBuilderExprString(expr))
		}

		currentObject = res.Val.(Object)
	}

	return currentObject, nil
}

func (p *Program) evaluateDotAccessAssignment(expr *gg_ast.DotAccessAssignmentExpression) error {
	target, err := p.getDotAccessAssignmentTarget(expr)
	if err != nil {
		return err
	}

	field := expr.Target.AccessChain[len(expr.Target.AccessChain)-1]
	val, err := p.evaluateValueExpr(expr.Value)
	if err != nil {
		return err
	}
	target[field] = val
	return nil
}
