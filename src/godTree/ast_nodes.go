package godTree

import "fmt"

// utility to prevent wrong assignments
type idKind int

const (
	IdExprNumber = idKind(ExprNumberLiteral)
	IdExprString = idKind(ExprStringLiteral)
	IdVariable   = idKind(ExprVariable)
)

type Identifier struct {
	Raw    string
	idKind idKind
}

func (id *Identifier) Name() string {
	return id.Raw
}

func (id *Identifier) Kind() ExpressionKind { return ExpressionKind(id.idKind) }

// a + b
type BinaryExpression struct {
	Lhs ValueExpression
	Op  string
	Rhs ValueExpression
}

func (be *BinaryExpression) Name() string         { return be.Op }
func (be *BinaryExpression) Kind() ExpressionKind { return ExprBinary }
func newBinaryExpression(lhs ValueExpression, operator string, rhs ValueExpression) *BinaryExpression {
	return &BinaryExpression{
		Lhs: lhs,
		Op:  operator,
		Rhs: rhs,
	}
}

// a 32
type AssignmentExpression struct {
	Target Identifier
	Value  ValueExpression
}

func (ae *AssignmentExpression) Kind() ExpressionKind { return ExprAssignment }
func newAssignmentExpression(target Identifier, value ValueExpression) *AssignmentExpression {
	return &AssignmentExpression{
		Target: target,
		Value:  value,
	}
}

func ind(count int) string {

	var spaces []rune
	spaces = append(spaces, '\n')
	for i := 0; i < count*2; i++ {
		spaces = append(spaces, ' ')
	}

	//var outRunes []rune
	//outRunes = append(outRunes, spaces...)
	//outRunes = append(outRunes, []rune(in)...)
	return string(spaces)
}

func ExprString(e Expression, depth int) string {
	switch e.Kind() {
	case ExprAssignment:
		val := e.(*AssignmentExpression)
		return fmt.Sprintf(ind(depth)+"assign of %v to %v", ExprString(val.Value, depth+1), val.Target)
	case ExprBinary:
		val := e.(*BinaryExpression)
		return fmt.Sprintf(ind(depth)+"operate %v %v %v", ExprString(val.Lhs, depth+1), val.Op, ExprString(val.Rhs, depth+1))
	case ExprNumberLiteral:
		goto IdentifierStr
	case ExprVariable:
		goto IdentifierStr
	case ExprStringLiteral:
		goto IdentifierStr
	default:
		panic(fmt.Sprintf("unknown expression type: %T", e))
	}

IdentifierStr:
	return fmt.Sprintf(ind(depth)+"Ident: %s", e.(*Identifier).Raw)
}
