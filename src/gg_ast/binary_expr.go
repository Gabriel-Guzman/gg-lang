package gg_ast

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	"gg-lang/src/operators"
	"gg-lang/src/token"
)

// iter should be pointing to right before the second expression after the first operator
// returns an identifier if there's no operator or a bin expression ready to be walked
// with operator precedence
func parseBinaryExpression(
	tokIter *iterator.Iter[token.Token],
) (IValExpr, error) {
	// begin building the initial binary expression if it exists
	lhsSve, err := parseSingleValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	op, ok := tokIter.Next()
	if !ok {
		return lhsSve, nil
	}
	if !op.TokenType.IsMathOperator() {
		return nil, ggErrs.Runtime("expected operator, got %s", tokIter.String())
	}

	rhs, err := parseSingleValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	lhs := &BinaryExpression{
		Lhs: lhsSve,
		Op:  op.Str,
		Rhs: rhs,
	}
	for {
		op, ok := tokIter.Next()
		if !ok {
			break
		}
		if !op.TokenType.IsMathOperator() {
			break
		}

		rhs, err := parseSingleValueExpr(tokIter)
		if err != nil {
			return nil, err
		}

		if operators.LeftFirst(lhs.Op, op.Str) {
			// left needs to be evaluated first and therefor deeper into the tree
			lhs = &BinaryExpression{
				Lhs: lhs,
				Op:  op.Str,
				Rhs: rhs,
			}
		} else {
			lhs.Rhs = &BinaryExpression{
				Lhs: lhs.Rhs,
				Op:  op.Str,
				Rhs: rhs,
			}
		}
	}

	return lhs, nil
}
