package godTree

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
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
	a.tokIter.Stringer = tokStringer
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
		if a.tokIter.HasNext() {
			return ggErrs.Runtime("couldnt finish parsing statement\n%s", a.tokIter.String())
		}
		// trap for function declaration
		if casted, ok := expr.(*FunctionDeclExpression); ok {
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

func newIdentifier(t tokenizer.Token) (*Identifier, error) {
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
	return &Identifier{Raw: t.Str, idKind: ik}, nil
}
