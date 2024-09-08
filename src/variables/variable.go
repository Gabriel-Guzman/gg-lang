//go:generate stringer -type=VarType
package variables

type VarType int

const (
	Integer VarType = iota
	String
	Boolean
	Function
	BuiltinFunction
	Void
)

type Variable struct {
	Name string
	//Typ   VarType
	Value *RuntimeValue
}
