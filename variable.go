//go:generate stringer -type=varType
package main

type varType int

const (
	INTEGER varType = iota
	STRING
)

type variable struct {
	name  string
	typ   varType
	value any
}
