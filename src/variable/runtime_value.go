package variable

import (
	"gg-lang/src/ggErrs"
	"strconv"
)

type RuntimeValue struct {
	Val interface{}
	Typ VarType
}

// if error is nil, return val is guaranteed to be of type targetType
func CoerceTo(val interface{}, targetType VarType) (interface{}, error) {
	switch val.(type) {
	case int:
		return CoerceFromInt(val.(int), targetType)
	case bool:
		return CoerceFromBool(val.(bool), targetType)
	default:
		return nil, ggErrs.Runtime("failed to coerce value of type %T to %s", val, targetType.String())
	}
}
func CoerceFromBool(val bool, targetType VarType) (interface{}, error) {
	switch targetType {
	case String:
		ret := strconv.FormatBool(val)
		return ret, nil
	case Boolean:
		return val, nil
	default:
		return nil, ggErrs.Runtime("failed to coerce bool %d to %s", val, targetType.String())
	}
}
func CoerceFromInt(val int, targetType VarType) (interface{}, error) {
	switch targetType {
	case String:
		ret := strconv.Itoa(val)
		return ret, nil
	case Integer:
		return val, nil
	default:
		return nil, ggErrs.Runtime("failed to coerce integer %d to %s", val, targetType.String())
	}
}
