package main

import (
	"context"
	"fmt"
	"os"

	"kivo/engine"
	"kivo/internal/cli"
)

func main() {
	ctx := context.Background()

	opts := engine.DefaultOpts()
	e, err := engine.New(ctx, opts)
	if err != nil {
		fmt.Println("error creating engine:", err)
		os.Exit(1)
	}

	repl := cli.New(e)
	repl.Run()
}
