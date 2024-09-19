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
		tokIter.Next()
		return nil, ggErrs.Runtime("unexpected token\n%s", tokIter.String())
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
	startingIndex := iter.Index()
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
		iter.SetIndex(startingIndex)
		binaryExpr, err := parseBinaryExpression(iter)
		if err != nil {
			return nil, err
		}
		// mark statement as done
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

func parseFuncDecl(tokIter *iterator.Iter[tokenizer.Token]) (*FunctionDeclExpression, error) {
	mbIdent, ok := tokIter.Next()
	if !ok {
		return nil, ggErrs.Runtime("Expected func name\n%s", tokIter.String())
	}

	id, err := newIdentifier(mbIdent)
	if err != nil {
		return nil, ggErrs.Runtime("Invalid identifier\n%s", tokIter.String())
	}

	mbOpenParen, ok := tokIter.Next()
	if !ok || mbOpenParen.TokenType != tokenizer.ROpenParen {
		return nil, ggErrs.Runtime("Expected (\n%s", tokIter.String())
	}

	var parms []string
	for {
		parm, ok := tokIter.Next()
		if !ok {
			return nil, ggErrs.Runtime("Unexpected end of param list\n%s", tokIter.String())
		}
		if parm.TokenType == tokenizer.RCloseParen {
			break
		}
		if parm.TokenType == tokenizer.RComma {
			continue
		}
		if parm.TokenType != tokenizer.Var {
			return nil, ggErrs.Runtime("Unexpected token\n%s", tokIter.String())
		}

		parms = append(parms, parm.Str)
	}

	mbOpenBrack, ok := tokIter.Next()
	if !ok || mbOpenBrack.TokenType != tokenizer.ROpenBrace {
		return nil, ggErrs.Runtime("Expected {\n%s", tokIter.String())
	}

	return &FunctionDeclExpression{
		Target: *id,
		Params: parms,
		Value:  nil,
	}, nil
}

// iter should be pointing to the opening parenthesis here
func newFuncCallExpression(funcName *Identifier, iter *iterator.Iter[tokenizer.Token]) (IValExpr, error) {
	nextTok, ok := iter.Peek()
	if !ok {
		return nil, ggErrs.Runtime("expected closing parenthesis ')' or args after function name\n%s", iter.String())
	}
	if nextTok.TokenType == tokenizer.RCloseParen {
		iter.Next() // consume the closing parenthesis ')'
		return &FunctionCallExpression{
			Id:   *funcName,
			Args: nil,
		}, nil
	}

	var args []IValExpr
	for {
		val, err := parseValueExpr(iter)
		if err != nil {
			return nil, err
		}
		if !iter.HasCurrent() {
			return nil, ggErrs.Runtime("unexpected end of arg list\n%s", iter.String())
		}
		args = append(args, val)
		mbComma := iter.Current()
		if mbComma.TokenType == tokenizer.RComma {
			continue
		}
		if mbComma.TokenType == tokenizer.RCloseParen {
			break
		}
	}
	return &FunctionCallExpression{
		Id:   *funcName,
		Args: args,
	}, nil
}
