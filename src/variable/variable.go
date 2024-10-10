//go:generate stringer -type=VarType
package variable

type VarType int

const (
	// Integer represents a whole number
	Integer VarType = iota
	// String represents a sequence of characters
	String
	// Boolean represents true and false values
	Boolean
	// Function is a user-defined function
	Function
	// BuiltinFunction is a function defined in the language's standard library
	BuiltinFunction
	// Object is a container of variables with a name
	Object
	// Void is the empty type, representing the absence of a value.
	Void
)

type Variable struct {
	Name         string
	RuntimeValue *RuntimeValue
}
