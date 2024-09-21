package gg_ast

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/parser"
	"gg-lang/src/token"
)

type astBuilder struct {
	StmtPar *parser.Parser[*parser.Parser[token.Token]]
}

func (a *astBuilder) parseBlockStatement() ([]Expression, error) {
	var expressions []Expression
	for a.StmtPar.HasCurr {
		tokParser := a.StmtPar.Curr
		if tokParser.Curr.TokenType == token.RCloseBrace {
			if tokParser.HasNext {
				return nil, ggErrs.Crit("unexpected token after closing brace - check tokenizer\n%s", tokParser.String())
			}

			a.StmtPar.Advance() // consume closing brace, which should be its own statement in this version
			break
		}

		stmt, err := parseTopLevelExpr(tokParser)
		if err != nil {
			return nil, err
		}
		a.StmtPar.Advance()
		expressions = append(expressions, stmt)
	}

	return nil, ggErrs.Syntax("no closing brace for block statement\n%s", a.StmtPar.String())
}

func newAstBuilder(ins [][]token.Token) *astBuilder {
	var stmtParserMembers []*parser.Parser[token.Token]
	for _, tokens := range ins {
		stmtParserMembers = append(stmtParserMembers, parser.New(tokens))
	}
	return &astBuilder{
		StmtPar: parser.New(stmtParserMembers),
	}
}

func BuildFromStatements(ins [][]token.Token) (*Ast, error) {
	a := newAstBuilder(ins)

	var expressions []Expression
	for a.StmtPar.HasCurr {
		stmt, err := parseTopLevelExpr(a.StmtPar.Curr)
		if err != nil {
			return nil, err
		}
		if !a.StmtPar.IsDone() {
			return nil, ggErrs.Crit("could not finish parsing statement\n%s", a.StmtPar.String())
		}

		// function trap, note continue at end of this block
		if stmt.Kind() == ExprFuncDecl {
			// function trap
			decl := stmt.(*FunctionDeclExpression)
			a.StmtPar.Advance() // consume the declaration
			exprs, err := a.parseBlockStatement()
			if err != nil {
				return nil, err
			}
			decl.Value = exprs
			continue
		}

		a.StmtPar.Advance()
	}

	return &Ast{Body: expressions}, nil
}
