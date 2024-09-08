package godTree

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	"gg-lang/src/operators"
	"gg-lang/src/tokenizer"
	"strings"
)

type Ast struct {
	Body     []Expression
	stmtIter *iterator.Iter[[]tokenizer.Token]
	tokIter  *iterator.Iter[tokenizer.Token]
}

func (a *Ast) nextStmt() bool {
	stmt, ok := a.stmtIter.Next()
	if !ok {
		return false
	}
	a.tokIter = iterator.New(stmt)
	return true
}

func (a *Ast) String() string {
	var sb strings.Builder
	for _, expr := range a.Body {
		ExprString(expr, 0, &sb)
	}
	return sb.String()
}

func New() *Ast {
	a := &Ast{}
	return a
}

func tokStringer(t tokenizer.Token) string {
	return t.Str
}

func (a *Ast) ParseStmts(tokens [][]tokenizer.Token) error {
	a.stmtIter = iterator.New(tokens)

	var currStmt []tokenizer.Token
	a.tokIter = iterator.New(currStmt)

	//nextStmt := func() bool {
	//	currStmt, ok := iter.Next()
	//	if !ok {
	//		return ok
	//	}
	//
	//	tokIter = iterator.New(currStmt)
	//	tokIter.Stringer = tokStringer
	//	return true
	//}

outer:
	for {
		ok := a.nextStmt()
		if !ok {
			break
		}
		expr, err := parseStmt(a.tokIter)
		if err != nil {
			return err
		}
		// trap for function declaration
		if casted, ok := expr.(*FunctionDeclExpression); ok {
			for {
				ok = a.nextStmt()
				if !ok {
					return ggErrs.Runtime("missing } in function decl\n%s", a.tokIter.String())
				}

				err := a.funcTrap(casted)
				if err != nil {
					return err
				}
				continue outer
			}
		}

		if a.tokIter.HasNext() {
			return ggErrs.Runtime("couldnt finish parsing statement\n%s", a.tokIter.String())
		}
		a.Body = append(a.Body, expr)
	}
	return nil
}

func (a *Ast) funcTrap(casted *FunctionDeclExpression) error {
	for {
		curr, ok := a.tokIter.Peek()
		if !ok {
			return ggErrs.Runtime("unexpected end of token iter in func trap\n%s", a.tokIter.String())
		}

		if curr.TokenType == tokenizer.RCloseBrace {
			a.Body = append(a.Body, casted)
			return nil
		}

		funcBodyExpr, err := parseStmt(a.tokIter)
		if err != nil {
			return err
		}
		casted.Value = append(casted.Value, funcBodyExpr)
		a.nextStmt()
	}
}

func parseStmt(tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	initial := tokIter.Copy()
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

	// should be an identifier
	id, err := newIdentifier(curr)
	if err != nil {
		return nil, err
	}

	next, ok := tokIter.Peek()
	// TODO this doesn't make sense
	if !ok {
		expr, err := parseValueExpr(tokIter)
		return expr, err
	}
	nextTokType := next.TokenType

	switch {
	case nextTokType == tokenizer.RAssign:
		tokIter.Next() // consume the '=' token
		expr, err := parseAssignmentExpr(id, tokIter)
		return expr, err
	case nextTokType.IsOperator():
		expr, err := parseValueExpr(initial)
		return expr, err
	default:
		return nil, ggErrs.Runtime("Unexpected token\n%s", tokIter.String)
	}
}

func parseFuncDecl(tokIter *iterator.Iter[tokenizer.Token]) (*FunctionDeclExpression, error) {
	mbIdent, ok := tokIter.Next()
	if !ok {
		return nil, ggErrs.Runtime("Expected func name\n%s", tokIter.String())
	}

	id, err := newIdentifier(mbIdent)
	if err != nil {
		return nil, err
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
		Parms:  parms,
		Value:  nil,
	}, nil
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
	tokIter := iter.Copy()
	// first word in value
	firstExpr, err := parseSingleValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	// tokIter is currently at the closing parenthesis ')'
	// go to tok after that
	afterParen, ok := tokIter.Next()
	if !ok {
		return firstExpr, nil
	}

	// operator, try for binary expression
	if afterParen.TokenType.IsOperator() {
		if afterParen.TokenType == tokenizer.RAssign {
			return nil, ggErrs.Runtime("invalid = in value expression: %s", tokIter.String())
		}

		binaryExpr, err := parseBinaryExpression(iter)
		if err != nil {
			return nil, err
		}
		return binaryExpr, nil
	}
	return firstExpr, nil
}

// iter should be pointing to right before the first tok in the value expression
func parseSingleValueExpr(tokIter *iterator.Iter[tokenizer.Token]) (IValExpr, error) {
	// first word in value
	firstTok, exists := tokIter.Next()

	if firstTok.TokenType.IsOperator() {
		return nil, ggErrs.Runtime("unexpected op %s at %d", firstTok.Str, firstTok.Start)
	}

	firstId, err := newIdentifier(firstTok)
	if err != nil {
		return nil, err
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

// iter should be pointing to right before the second expression after the first operator
// returns a bin expression ready to be walked with operator precedence
func parseBinaryExpression(
	tokIter *iterator.Iter[tokenizer.Token],
) (*BinaryExpression, error) {
	secondExpr, err := parseSingleValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	firstBinaryExpr := &BinaryExpression{
		Lhs: firstSingleValueExpr,
		Op:  op,
		Rhs: secondExpr,
	}

	mbOp, ok := tokIter.Peek()
	if !ok {
		return firstBinaryExpr, nil
	}
	if !mbOp.TokenType.IsOperator() {
		return nil, ggErrs.Runtime("unexpected token %s", tokIter.String())
	}

	tokIter.Next() // consume the operator token

	secondBinExpr, err := parseBinaryExpression(secondExpr, mbOp.Str, tokIter)
	if err != nil {
		return nil, err
	}

	// note firstBinaryExpr.Rhs == secondBinExpr.Lhs
	if operators.LeftFirst(firstBinaryExpr.Op, secondBinExpr.Op) {
		// firstBinaryExpr should go before mbOp, so deeper in the tree than it
		return &BinaryExpression{
			Lhs: firstBinaryExpr,
			Op:  secondBinExpr.Op,
			Rhs: secondBinExpr.Rhs,
		}, nil
	} else {
		//firstBinaryExpr should go after mbOp, so shallower in the tree than it
		return &BinaryExpression{
			Lhs: firstBinaryExpr.Lhs,
			Op:  firstBinaryExpr.Op,
			Rhs: secondBinExpr,
		}, nil
	}
}

func newIdentifier(t tokenizer.Token) (*Identifier, error) {
	var ik IdExprKind
	switch t.TokenType {
	case tokenizer.NumberLiteral:
		ik = IdExprNumber
	case tokenizer.Var:
		ik = IdExprVariable
	case tokenizer.StringLiteral:
		ik = IdExprString
	case tokenizer.TrueLiteral:
		ik = IdExprBool
	case tokenizer.FalseLiteral:
		ik = IdExprBool
	default:
		return nil, ggErrs.Runtime("invalid identifier %s at %d: ", t.Str, t.Start)
	}
	return &Identifier{Raw: t.Str, idKind: ik}, nil
}
