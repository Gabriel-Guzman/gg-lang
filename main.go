package main

import (
	"gg-lang/src/schemes"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		schemes.Repl()
		return
	}

	if len(os.Args) != 2 {
		panic("Usage: go run main.go <filename>")
	}

	// get arguments
	filename := os.Args[1]
	schemes.Exec(filename)
}
