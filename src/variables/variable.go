//go:generate stringer -type=varType
package variables

type VarType int

const (
	INTEGER VarType = iota
	STRING
)

type Variable struct {
	Name  string
	Typ   VarType
	Value any
}
