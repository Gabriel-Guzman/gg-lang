package gg_ast

import (
	"gg-lang/src/gg"
	"gg-lang/src/operators"
	"gg-lang/src/parser"
	"gg-lang/src/token"
)

/*
This is the top-level expression parser.
Its job is to
 1. pick the appropriate expression parser (all defined below) based on the
		first two tokens in the parser argument.
 2. parse the expression

The parser should be at its initial state, i.e. with the current index set to 0 and
pointing to the first token in the expression.
After a successful parse, the parser should be pointing to the token after the expression
*/

type tokenParser = *parser.Parser[token.Token]

// a statement is a function call, declaration, assignment expression, or a for loop expression
// it is up to the builder to disallow these expressions if it's not parsing the top level
func parseExpression(p tokenParser) (Expression, error) {
	// first expression should be an identifier, a reserved keyword, or a value expression
	// check reserved keywords first
	if p.Curr.TokenType == token.Function {
		return parseFuncDecl(p)
	}
	if p.Curr.TokenType == token.Try {
		return parseTryCatchExpr(p)
	}
	if p.Curr.TokenType == token.For {
		return parseForLoopExpr(p)
	}
	if p.Curr.TokenType == token.If {
		return parseIfElseExpr(p)
	}
	if p.Curr.TokenType == token.Return {
		return parseReturnExpr(p)
	}
	if p.Curr.TokenType == token.OpenBrace {
		return parseObjectExpr(p)
	}
	if p.Curr.TokenType == token.OpenParen {
		return parseParenExpr(p)
	}
	if p.Curr.TokenType == token.OpenBracket {
		return parseArrayDeclExpr(p)
	}

	// now it could be a function call or an assignment expression, both of which
	// have to start with an identifier. this means no unassigned value expressions
	// other than function calls are allowed at the top level.
	id, err := parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if p.Curr.TokenType == token.Dot {
		expr, err := parseDotAccessExpr(id, p)
		if err != nil {
			return nil, err
		}
		if advanceIfCurrIs(p, token.Assign) {
			return parseDotAccessAssignExpr(expr, p)
		}
	}

	if p.Curr.TokenType == token.Assign {
		p.Advance()
		return parseAssignmentExpr(id, p)
	}

	if p.Curr.TokenType == token.OpenParen {
		expr, err := parseFuncCallExpr(id, p)
		if err != nil {
			return nil, err
		}
		if !advanceIfCurrIs(p, token.Term) {
			return nil, gg.Syntax("expected ; after top-level function call\n%s", p.String())
		}
		return expr, nil
	}

	if p.Curr.TokenType == token.OpenBracket {
		expr, err := parseArrayAccessExpr(id, p)
		if err != nil {
			return nil, err
		}
		return parseArrayIndexAssignExpr(expr, p)
	}

	return nil, gg.Syntax("invalid top-level expression: %s \nin %s", p.Curr.Symbol, p.String())

	////// if no operator, it's a function call
	//expr, err := parseFuncCallExpr(id, p)
	//if err != nil {
	//	return nil, err
	//}
	//if !advanceIfCurrIs(p, token.Term) {
	//	return nil, gg.Syntax("expected ; after function call\n%s", p.String())
	//}
	//return expr, nil
}

/*
The parser argument to an expression parser should
be pointing to the first token in the expression to be parsed,
meaning parser.Curr == <first relevant token>
If the expression starts with a keyword, p will be pointing to the next token.
After a successful parse, the parser should be pointing to the token after the expression
*/
func parseObjectExpr(p tokenParser) (ValueExpression, error) {
	if !advanceIfCurrIs(p, token.OpenBrace) {
		return nil, gg.Syntax("expected opening brace for object expression\n%s", p.String())
	}
	props := make(map[string]ValueExpression)
	for p.HasCurr && p.Curr.TokenType != token.CloseBrace {
		prop, err := parseIdentifier(p)
		if err != nil {
			return nil, err
		}
		if prop.Kind() != ExprVariable {
			return nil, gg.Crit("expected identifier as object property name, got %s instead in\n%s", prop.Name(), p.String())
		}
		if !advanceIfCurrIs(p, token.Colon) {
			return nil, gg.Syntax("expected ':' after object property name\n%s", p.String())
		}
		expr, err := parseValueExpr(p)

		if err != nil {
			return nil, err
		}
		props[prop.Tok.Symbol] = expr
		if !advanceIfCurrIs(p, token.Comma) {
			break
		}
	}
	if !advanceIfCurrIs(p, token.CloseBrace) {
		return nil, gg.Syntax("expected closing brace for object expression\n%s", p.String())
	}
	return &ObjectExpression{Properties: props}, nil
}

func parseTryCatchExpr(p tokenParser) (Expression, error) {
	if !advanceIfCurrIs(p, token.Try) {
		return nil, gg.Syntax("expected 'try' keyword for try-catch expression\n%s", p.String())
	}
	tryBlock, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	if !advanceIfCurrIs(p, token.Catch) {
		return nil, gg.Syntax("expected 'catch' keyword for try-catch expression\n%s", p.String())
	}

	parenParams, err := params(p, token.OpenParen, token.CloseParen)
	if err != nil {
		return nil, err
	}
	if len(parenParams) != 1 {
		return nil, gg.Syntax("catch statement requires 1 parameter\n%s", p.String())
	}
	catchBlock, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	expr := &TryCatchExpression{
		Try: &tryBlock,
		Catch: &CatchExpression{
			ErrorParam: parenParams[0].Symbol,
			Body:       &catchBlock,
		},
	}

	if advanceIfCurrIs(p, token.Finally) {
		finallyBlock, err := parseBlockStatement(p)
		if err != nil {
			return nil, err
		}
		expr.Finally = &finallyBlock
	}

	return expr, nil
}

func parseParenExpr(p tokenParser) (ValueExpression, error) {
	if !advanceIfCurrIs(p, token.OpenParen) {
		return nil, gg.Syntax("expected opening parenthesis for parenthesized expression\n%s", p.String())
	}
	expr, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}
	if !advanceIfCurrIs(p, token.CloseParen) {
		return nil, gg.Syntax("expected closing parenthesis for parenthesized expression\n%s", p.String())
	}
	return &ParenthesizedExpression{Expr: expr}, nil
}

func parseDotAccessExpr(id *Identifier, p tokenParser) (*DotAccessExpression, error) {
	if !advanceIfCurrIs(p, token.Dot) {
		return nil, gg.Syntax("expected '.' after dot access expression\n%s", p.String())
	}

	chain := []string{id.Name()}
	for p.HasCurr && p.Curr.TokenType == token.Ident {
		chain = append(chain, p.Curr.Symbol)
		p.Advance()
		if !advanceIfCurrIs(p, token.Dot) {
			break
		}
	}

	if len(chain) <= 1 {
		return nil, gg.Syntax("expected identifier after '.'\n%s", p.String())
	}

	return &DotAccessExpression{AccessChain: chain}, nil
}

func parseAssignmentExpr(target *Identifier, p tokenParser) (*AssignmentExpression, error) {
	expr, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}
	if !advanceIfCurrIs(p, token.Term) {
		return nil, gg.Syntax("expected ; after assignment expression\n%s", p.String())
	}

	return &AssignmentExpression{Target: target, Value: expr}, nil
}

func parseForLoopExpr(p tokenParser) (*ForLoopExpression, error) {
	if !advanceIfCurrIs(p, token.For) { // eat the for keyword
		return nil, gg.Crit("expected 'for' keyword in expression parser\n%s", p.String())
	}

	condition, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}

	body, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}
	return &ForLoopExpression{Condition: condition, Body: body}, nil
}

func parseIfElseExpr(p tokenParser) (*IfElseStatement, error) {
	if !advanceIfCurrIs(p, token.If) { // eat the if keyword
		return nil, gg.Crit("expected 'if' keyword in expression parser\n%s", p.String())
	}
	res := &IfElseStatement{}
	condition, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}
	res.Condition = condition

	body, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}
	res.Body = body

	if advanceIfCurrIs(p, token.Else) {
		if p.Curr.TokenType == token.If {
			alt, err := parseIfElseExpr(p)
			if err != nil {
				return nil, err
			}
			res.ElseExpression = alt
		} else {
			b, err := parseBlockStatement(p)
			if err != nil {
				return nil, err
			}
			res.ElseExpression = b
		}
	}

	return res, nil
}

func parseReturnExpr(p tokenParser) (*ReturnStatement, error) {
	if !advanceIfCurrIs(p, token.Return) { // eat the return keyword
		return nil, gg.Crit("expected 'return' keyword in expression parser\n%s", p.String())
	}

	expr, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}

	if !advanceIfCurrIs(p, token.Term) {
		return nil, gg.Syntax("expected ; after return expression\n%s", p.String())
	}

	return &ReturnStatement{Value: expr}, nil
}

func params(p tokenParser, open token.Type, close token.Type) ([]token.Token, error) {
	if !advanceIfCurrIs(p, open) {
		return nil, gg.Runtime("expected '(' after function name\n%s", p.String())
	}

	var params []token.Token
	for {
		if !p.HasCurr {
			return nil, gg.Runtime("Unexpected end of param list\n%s", p.String())
		}
		param := p.Curr
		if param.TokenType == close {
			p.Advance() // eat the closing parenthesis ')'
			break
		}
		if param.TokenType == token.Comma {
			p.Advance() // eat the comma
			continue
		}
		if param.TokenType != token.Ident {
			return nil, gg.Runtime("Unexpected token\n%s", p.String())
		}

		params = append(params, param)
		p.Advance()
	}
	return params, nil
}

func parseFuncDecl(p tokenParser) (*FunctionDeclExpression, error) {
	p.Advance() // eat the function keyword

	id, err := parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	params, err := params(p, token.OpenParen, token.CloseParen)
	if err != nil {
		return nil, err
	}

	block, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	return &FunctionDeclExpression{
		Target: id,
		Params: params,
		Body:   block,
	}, nil
}

func parseDotAccessAssignExpr(target *DotAccessExpression, p tokenParser) (*DotAccessAssignmentExpression, error) {
	val, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}

	if !advanceIfCurrIs(p, token.Term) {
		return nil, gg.Syntax("expected ; after dot-access assignment expression\n%s", p.String())
	}

	return &DotAccessAssignmentExpression{
		Target: target,
		Value:  val,
	}, nil
}

func parseArrayDeclExpr(p tokenParser) (*ArrayDeclExpression, error) {
	members, err := arguments(p, token.OpenBracket, token.CloseBracket)
	if err != nil {
		return nil, err
	}
	return &ArrayDeclExpression{Elements: members}, nil
}

func parseIdentifier(p tokenParser) (*Identifier, error) {
	t := p.Curr
	var ik IdExprKind
	switch t.TokenType {
	case token.IntLiteral:
		ik = IdExprNumber
	case token.Ident:
		ik = IdExprVariable
	case token.StringLiteral:
		ik = IdExprString
	case token.TrueLiteral:
		ik = IdExprBool
	case token.FalseLiteral:
		ik = IdExprBool
	default:
		return nil, gg.Runtime("invalid identifier %s, in\n$s", t.Symbol, p.String())
	}
	p.Advance()
	return &Identifier{Tok: t, idKind: ik}, nil
}

// returns a primary expression or a binary expression
func parseValueExpr(p tokenParser) (ValueExpression, error) {
	if !p.HasCurr {
		return nil, gg.Runtime("unexpected end of expression\n%s", p.String())
	}

	// build initial binary tree
	lhsNonBinary, err := parsePrimaryExpr(p)
	if err != nil {
		return nil, err
	}

	if !p.Curr.TokenType.IsOperator() {
		return lhsNonBinary, nil
	}

	op := p.Curr
	p.Advance() // eat the operator token

	rhs, err := parsePrimaryExpr(p)
	if err != nil {
		return nil, err
	}

	lhs := &BinaryExpression{
		Lhs: lhsNonBinary,
		Op:  op,
		Rhs: rhs,
	}

	// add on to the initial tree
	for {
		op = p.Curr
		if !op.TokenType.IsOperator() {
			break
		}
		p.Advance() // eat the operator token

		rhs, err := parsePrimaryExpr(p)
		if err != nil {
			return nil, err
		}

		if operators.LeftFirst(lhs.Op.Symbol, op.Symbol) {
			// left needs to be evaluated first and therefor deeper into the tree
			lhs = &BinaryExpression{
				Lhs: lhs,
				Op:  op,
				Rhs: rhs,
			}
		} else {
			lhs.Rhs = &BinaryExpression{
				Lhs: lhs.Rhs,
				Op:  op,
				Rhs: rhs,
			}
		}
	}
	return lhs, nil
}

// A primary expression is either an identifier, a literal, a function call, a unary binary expression, or a function declaration
func parsePrimaryExpr(p tokenParser) (ValueExpression, error) {
	if !p.HasCurr {
		return nil, gg.Runtime("unexpected end of expression\n%s", p.String())
	}

	// unary operators
	if p.Curr.TokenType.IsOperator() {
		op := p.Curr
		p.Advance()
		toNegate, err := parsePrimaryExpr(p)
		if err != nil {
			return nil, err
		}

		return &UnaryExpression{
			Op:  op,
			Rhs: toNegate,
		}, nil
	}

	if p.Curr.TokenType == token.Function {
		return parseFuncDecl(p)
	}

	if p.Curr.TokenType == token.OpenBrace {
		return parseObjectExpr(p)
	}

	if p.Curr.TokenType == token.OpenParen {
		return parseParenExpr(p)
	}

	if p.Curr.TokenType == token.OpenBracket {
		return parseArrayDeclExpr(p)
	}

	id, err := parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if p.Curr.TokenType == token.Dot {
		return parseDotAccessExpr(id, p)
	}

	if p.Curr.TokenType == token.OpenParen {
		return parseFuncCallExpr(id, p)
	}

	if p.Curr.TokenType == token.OpenBracket {
		return parseArrayAccessExpr(id, p)
	}

	return id, nil
}

func parseArrayAccessExpr(id *Identifier, p tokenParser) (*ArrayIndexExpression, error) {
	args, err := arguments(p, token.OpenBracket, token.CloseBracket)
	if err != nil {
		return nil, err
	}
	if len(args) != 1 {
		return nil, gg.Syntax("expected 1 argument in array access expression, got %d\n%s", len(args), p.String())
	}
	return &ArrayIndexExpression{
		Array: id,
		Index: args[0],
	}, nil
}

func parseArrayIndexAssignExpr(arrayAccessExpr *ArrayIndexExpression, p tokenParser) (*ArrayIndexAssignmentExpression, error) {
	if !advanceIfCurrIs(p, token.Assign) {
		return nil, gg.Syntax("expected '=' after array index assignment expression\n%s", p.String())
	}

	val, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}
	if !advanceIfCurrIs(p, token.Term) {
		return nil, gg.Syntax("expected ; after array index assignment\n%s", p.String())
	}

	return &ArrayIndexAssignmentExpression{
		ArrayIndexExpression: arrayAccessExpr,
		Value:                val,
	}, nil
}

func arguments(p tokenParser, open, close token.Type) ([]ValueExpression, error) {
	if !advanceIfCurrIs(p, open) {
		return nil, gg.Syntax("expected '%s' after '%s'\n%s", close, open, p.String())
	}

	var args []ValueExpression
	if advanceIfCurrIs(p, close) {
		return args, nil
	}

	for {
		if p.HasCurr {
			expr, err := parseValueExpr(p)
			if err != nil {
				return nil, err
			}
			args = append(args, expr)

			if !advanceIfCurrIs(p, token.Comma) {
				break
			}
		} else {
			return nil, gg.Syntax("unexpected end of expression\n%s", p.String())
		}
	}

	if !advanceIfCurrIs(p, close) {
		return nil, gg.Syntax("expected '%s' after last argument\n%s", close, p.String())
	}

	return args, nil
}
func parseFuncCallExpr(id *Identifier, p tokenParser) (ValueExpression, error) {
	args, err := arguments(p, token.OpenParen, token.CloseParen)
	if err != nil {
		return nil, err
	}

	return &FunctionCallExpression{
		Id:   id,
		Args: args,
	}, nil
}

// Advances the parser if the current token matches the given token type
func advanceIfCurrIs(p tokenParser, tt token.Type) bool {
	return p.AdvanceIf(func(t token.Token) bool { return t.TokenType == tt })
}
