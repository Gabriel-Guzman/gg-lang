//go:generate stringer -type=token
package main

type tokenType int

const (
	VAR tokenType = iota
	OPERATOR
	NUMBER_LITERAL
	STRING_LITERAL
)

type token struct {
	start     int
	end       int
	str       string
	tokenType tokenType
}
