package program

import (
	"gg-lang/src/gg"
	"gg-lang/src/gg_ast"
	"gg-lang/src/variable"
)

type Array = []variable.RuntimeValue

func (p *Program) evaluateArrayDeclExpression(expr *gg_ast.ArrayDeclExpression) (*variable.RuntimeValue, error) {
	val := make(Array, len(expr.Elements))
	for i, elem := range expr.Elements {
		v, err := p.evaluateValueExpr(elem)
		if err != nil {
			return nil, err
		}
		val[i] = *v
	}

	return &variable.RuntimeValue{
		Val: val,
		Typ: variable.Array,
	}, nil
}

func (p *Program) evaluateArrayIndexExpression(expr *gg_ast.ArrayIndexExpression) (*variable.RuntimeValue, error) {
	index, err := p.evaluateValueExpr(expr.Index)
	if err != nil {
		return nil, err
	}
	arr := p.currentScope().findVariable(expr.Array.Name())
	arrVal, ok := arr.RuntimeValue.Val.(Array)
	if !ok {
		return nil, gg.Runtime("array index expression must reference an array\n%+v", expr)
	}

	var indexVal int
	if val, ok := index.Val.(int); !ok {
		return nil, gg.Runtime("array index must evaluate to int\n%+v", expr)
	} else {
		indexVal = val
	}

	length := len(arrVal)
	if indexVal < 0 || indexVal >= length {
		return nil, gg.Runtime("array index out of range\n%+v", expr)
	}

	return &variable.RuntimeValue{
		Val: arrVal[indexVal].Val,
		Typ: arrVal[indexVal].Typ,
	}, nil
}

func (p *Program) evaluateArrayIndexAssignmentExpression(expr *gg_ast.ArrayIndexAssignmentExpression) error {
	arr := p.currentScope().findVariable(expr.Array.Name())
	if arr == nil {
		return gg.Runtime("undefined variable: %s\n%+v", expr.Array.Name(), expr)
	}

	// Check if array is actually an array
	arrVal, ok := arr.RuntimeValue.Val.(Array)
	if !ok {
		return gg.Runtime("array index assignment expression must reference an array\n%+v", expr)
	}

	// Check if index is an integer within bounds
	indexVal, err := p.evaluateValueExpr(expr.Index)
	if err != nil {
		return err
	}

	index, ok := indexVal.Val.(int)
	if !ok {
		return gg.Runtime("array index must evaluate to int\n%+v", expr)
	}

	if index < 0 || index >= len(arrVal) {
		return gg.Runtime("array index out of range\n%+v", expr)
	}

	// evaluate right side of the assignment expression
	newVal, err := p.evaluateValueExpr(expr.Value)
	if err != nil {
		return err
	}

	arrVal[index] = *newVal
	return nil
}
