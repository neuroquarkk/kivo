package main

import (
	"context"
	"fmt"
	"kivo/engine"
	"kivo/internal/cli"
	"os"
)

func main() {
	ctx := context.Background()

	opts := &engine.Opts{}
	e, err := engine.New(ctx, opts)
	if err != nil {
		fmt.Println("error creating engine:", err)
		os.Exit(1)
	}

	repl := cli.New(e)
	repl.Run()
}
