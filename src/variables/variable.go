//go:generate stringer -type=VarType
package variables

type VarType int

const (
	Integer VarType = iota
	String
	Function
)

type Variable struct {
	Name  string
	Typ   VarType
	Value any
}
