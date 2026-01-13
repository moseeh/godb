package main

import (
	"godb/repl"
	"os"
)

func main() {
	r := repl.NewREPL(os.Stdin)
	r.Start()
}
