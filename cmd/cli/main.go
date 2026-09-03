package main

import (
	"context"
	"kivo/engine"
	"kivo/internal/cli"
)

func main() {
	ctx := context.Background()

	e := engine.New(ctx)
	repl := cli.New(e)
	repl.Run()
}
