package godTree

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/parser"
	"gg-lang/src/tokenizer"
)

/*
This is the top-level expression parser.
Its job is to
 1. pick the appropriate expression parser (all defined below) based on

the first two tokens in the parser argument.
 2. parse the expression
 3. handle block statements (if any)

The parser should be at its initial state, i.e. with the current index set to 0 and
pointing to the first token in the expression.
After a successful parse, the parser should be pointing to the token after the expression
*/
func _parseStatement(p *parser.Parser[tokenizer.Token]) (Expression, error) {
	// first expression should be an identifier, a reserved keyword, or a value expression

	// check reserved keywords first
	if p.Curr.TokenType == tokenizer.Function {
		funcDecl, err := _parseFuncDecl(p)
		if err != nil {
			return nil, err
		}
		return funcDecl, nil
	}

	// now it could be a function call or an assignment expression, both of which
	// have to start with an identifier. this means no unassigned value expressions
	// other than function calls are allowed.
	id, err := _parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if p.Curr.TokenType == tokenizer.RAssign {
		p.Advance()
		return _parseAssignmentExpr(id, p)
	}

	// if no operator, it's a function call
	return _parseFuncCallExpr(id, p)
}

/*
The parser argument to an expression parser should
be pointing to the first token in the expression to be parsed,
meaning parser.Curr == <first relevant token>
This will almost always be the token right after the first identifier in the expression
i.e. the "(" in "funCall();" or the "=" in "x = 1 + 2;"
After a successful parse, the parser should be pointing to the token after the expression
*/
func _parseAssignmentExpr(target *Identifier, p *parser.Parser[tokenizer.Token]) (*AssignmentExpression, error) {
	expr, err := _parseValueExpr(p)
	if err != nil {
		return nil, err
	}

	return &AssignmentExpression{Target: target, Value: expr}, nil
}

func _parseFuncDecl(p *parser.Parser[tokenizer.Token]) (*FunctionDeclExpression, error) {
	p.Advance() // eat the function keyword

	id, err := _parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if !advanceIfCurrIs(p, tokenizer.ROpenParen) {
		return nil, ggErrs.Runtime("expected '(' after function name\n%s", p.String())
	}

	var params []string
	for {
		if !p.HasCurr {
			return nil, ggErrs.Runtime("Unexpected end of param list\n%s", p.String())
		}
		param := p.Curr
		if param.TokenType == tokenizer.RCloseParen {
			break
		}
		if param.TokenType == tokenizer.RComma {
			continue
		}
		if param.TokenType != tokenizer.Ident {
			return nil, ggErrs.Runtime("Unexpected token\n%s", p.String())
		}

		params = append(params, param.Str)
		p.Advance()
	}

	if !advanceIfCurrIs(p, tokenizer.ROpenBrace) {
		return nil, ggErrs.Runtime("expected '{' after function declaration\n%s", p.String())
	}

	return &FunctionDeclExpression{
		Target: id,
		Params: params,
	}, nil
}

func _parseIdentifier(p *parser.Parser[tokenizer.Token]) (*Identifier, error) {
	t := p.Curr
	var ik IdExprKind
	switch t.TokenType {
	case tokenizer.IntLiteral:
		ik = IdExprNumber
	case tokenizer.Ident:
		ik = IdExprVariable
	case tokenizer.StringLiteral:
		ik = IdExprString
	case tokenizer.TrueLiteral:
		ik = IdExprBool
	case tokenizer.FalseLiteral:
		ik = IdExprBool
	default:
		return nil, ggErrs.Runtime("invalid identifier %s", t.Str)
	}
	p.Advance()
	return &Identifier{Raw: t.Str, idKind: ik}, nil
}

func _parseValueExpr(p *parser.Parser[tokenizer.Token]) (IValExpr, error) {
	if !p.HasCurr {
		return nil, ggErrs.Runtime("unexpected end of expression\n%s", p.String())
	}

	// unary operators
	if p.Curr.TokenType == tokenizer.RMinus {
		//lhs := newBinaryExpression()
	}

	switch {
	case p.Curr.TokenType.IsOperator():
		return _parseBinaryExpr(p)
	}
}

func _parseBinaryExpr(p *parser.Parser[tokenizer.Token]) (IValExpr, error) {

}

func _parseSingleValueExpr(id *Identifier, p *parser.Parser[tokenizer.Token]) (IValExpr, error) {

	return nil, ggErrs.Crit("not implemented")
}

func _parseFuncCallExpr(id *Identifier, p *parser.Parser[tokenizer.Token]) (IValExpr, error) {
	if !advanceIfCurrIs(p, tokenizer.ROpenParen) {
		return nil, ggErrs.Runtime("expected '(' after function name\n%s", p.String())
	}

	var args []IValExpr
	for {
		if !p.HasCurr {
			return nil, ggErrs.Runtime("Unexpected end of arg list\n%s", p.String())
		}
		arg := p.Curr
		if arg.TokenType == tokenizer.RCloseParen {
			break
		}
		if arg.TokenType == tokenizer.RComma {
			continue
		}
		expr, err := _parseValueExpr(p)
		if err != nil {
			return nil, err
		}

		args = append(args, expr)
		p.Advance()
	}

	return &FunctionCallExpression{
		Id:   id,
		Args: args,
	}, nil
}

// Advances the parser if the current token matches the given token type
func advanceIfCurrIs(p *parser.Parser[tokenizer.Token], tt tokenizer.TokenType) bool {
	if p.HasCurr && p.Curr.TokenType == tt {

		p.Advance()
		return true
	}

	return false
}
