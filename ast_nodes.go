package main

type identifier struct {
	raw string
}

func (id *identifier) name() string {
	return id.raw
}

func (id *identifier) kind() expressionKind { return ExprIdentifier }

func newIdentifier(w word) *identifier {
	return &identifier{raw: w.str}
}

type numberLiteral struct {
	raw string
}

func (nl *numberLiteral) name() string {
	return nl.raw
}

func (nl *numberLiteral) kind() expressionKind { return ExprNumberLiteral }

//func (nl *numberLiteral) name() (interface{}, error) {
//	num, err := strconv.Atoi(nl.raw)
//	if err != nil {
//		return nil, fmt.Errorf("invalid number literal: %s", nl.raw)
//	}
//	return num, nil
//}

// a + b
type binaryExpression struct {
	lhs      valueExpression
	operator word
	rhs      expression
}

func (be *binaryExpression) kind() expressionKind { return ExprBinary }
func newBinaryExpression(lhs valueExpression, operator word, rhs expression) *binaryExpression {
	return &binaryExpression{
		lhs:      lhs,
		operator: operator,
		rhs:      rhs,
	}
}

// a 32
type assignmentExpression struct {
	target identifier
	value  expression
}

func (ae *assignmentExpression) kind() expressionKind { return ExprAssignment }
func newAssignmentExpression(target identifier, value expression) *assignmentExpression {
	return &assignmentExpression{
		target: target,
		value:  value,
	}
}
