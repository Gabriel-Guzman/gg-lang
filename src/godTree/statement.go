package godTree

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	"gg-lang/src/tokenizer"
)

func parseStmt(tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	// move to first token in stmt
	curr, exists := tokIter.Next()
	if !exists {
		return nil, ggErrs.Runtime("expected a statement\n%s", tokIter.String())
	}

	// check for reserved keywords
	if curr.TokenType == tokenizer.Function {
		funcDecl, err := parseFuncDecl(tokIter)
		return funcDecl, err
	}

	// should be a single value expr
	tokIter.Reset()
	firstSingleValueExpr, err := parseSingleValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	next, ok := tokIter.Peek()
	if !ok {
		return firstSingleValueExpr, err
	}
	nextTokType := next.TokenType

	switch {
	case nextTokType == tokenizer.RAssign:
		tokIter.Next() // consume the '=' token
		id, ok := firstSingleValueExpr.(*Identifier)
		if !ok {
			return nil, ggErrs.Runtime("Expected identifier before '=', got %s", tokIter.String())
		}
		expr, err := parseAssignmentExpr(id, tokIter)
		return expr, err
	default:
		tokIter.Reset()
		return parseValueExpr(tokIter)
	}
}

// tokIter should be pointing to the token right before the value expression
func parseAssignmentExpr(id *Identifier, tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	expr, err := parseValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	at := newAssignmentExpression(id, expr)
	return at, nil
}

// tokIter should be pointing to the token right before the value expression
func parseValueExpr(iter *iterator.Iter[tokenizer.Token]) (IValExpr, error) {
	beginning := iter.Copy()
	// first word in value
	firstExpr, err := parseSingleValueExpr(iter)
	if err != nil {
		return nil, err
	}

	// iter is currently at the closing parenthesis ')'
	// go to tok after that
	afterParen, ok := iter.Next()
	if !ok {
		return firstExpr, nil
	}

	// operator, try for binary expression
	if afterParen.TokenType.IsOperator() {
		if afterParen.TokenType == tokenizer.RAssign {
			return nil, ggErrs.Runtime("invalid = in value expression: %s", iter.String())
		}

		// point to the first identifier in the binary expression
		iter.Reset()
		binaryExpr, err := parseBinaryExpression(beginning)
		if err != nil {
			return nil, err
		}
		// mark statement as done
		iter.End()
		return binaryExpr, nil
	}
	return firstExpr, nil
}

// iter should be pointing to right before the first tok in the value expression
func parseSingleValueExpr(tokIter *iterator.Iter[tokenizer.Token]) (IValExpr, error) {
	// first word in value
	firstTok, exists := tokIter.Next()

	if firstTok.TokenType.IsOperator() {
		if firstTok.TokenType == tokenizer.RMinus {
			// unary minus
			next, err := parseSingleValueExpr(tokIter)
			if err != nil {
				return nil, err
			}

			return &BinaryExpression{
				Lhs: &Identifier{
					Raw:    "-1",
					idKind: IdExprNumber,
				},
				Op:  "*",
				Rhs: next,
			}, nil
		}
		return nil, ggErrs.Runtime("unexpected op %s", tokIter.String())
	}

	firstId, err := newIdentifier(firstTok)
	if err != nil {
		return nil, ggErrs.Runtime("could not parse identifier\n%s", tokIter.String())
	}

	// try for next word
	nextTok, exists := tokIter.Peek()
	if !exists {
		// no next word, return identifier
		return firstId, nil
	}

	var firstExpr IValExpr
	firstExpr = firstId
	// try for function call
	if nextTok.TokenType == tokenizer.ROpenParen {
		funcName := firstId
		// point to the opening parenthesis
		tokIter.Next()
		firstExpr, err = newFuncCallExpression(funcName, tokIter)
		if err != nil {
			return nil, err
		}
	}

	return firstExpr, nil
}
