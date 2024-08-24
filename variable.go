package main

type varType int

const (
	INTEGER varType = iota
)

type variable struct {
	name  string
	typ   varType
	value interface{}
}
