package main

import (
	"kivo/engine"
	"kivo/internal/cli"
)

func main() {
	e := engine.New()
	repl := cli.New(e)
	repl.Run()
}
