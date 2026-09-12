package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/reeflective/readline"
	"kivo/engine"
)

type REPL struct {
	engine *engine.Engine
}

func New(e *engine.Engine) *REPL {
	return &REPL{engine: e}
}

func (r *REPL) Run() {
	rl := readline.NewShell()
	rl.Prompt.Primary(func() string { return "kivo> " })

	for {
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				continue
			}
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Println("input error:", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if shouldExit := r.dispatch(line); shouldExit {
			return
		}
	}
}

func (r *REPL) dispatch(line string) (exit bool) {
	args := parseArgs(line)
	if len(args) == 0 {
		return false
	}

	cmd, rest := args[0], args[1:]

	switch strings.ToUpper(cmd) {
	case "SET":
		r.cmdSet(rest)
	case "GET":
		r.cmdGet(rest)
	case "DELETE":
		r.cmdDelete(rest)
	case "EXISTS":
		r.cmdExists(rest)
	case "COUNT":
		r.cmdCount()
	case "INFO":
		r.cmdInfo()
	case "FLUSH":
		r.cmdFlush()
	case "HELP":
		printHelp()
	case "EXIT", "QUIT":
		return true
	default:
		fmt.Printf("unknown command: %s\n", cmd)
	}

	return false
}

func printHelp() {
	fmt.Println(`commands:
  SET <key> <value> [ttl]   store a value, optional TTL (30s, 5m, 1h)
  GET <key>                 retrieve a value
  DELETE <key>           	remove a key
  EXISTS <key>              check whether a key exists
  COUNT                     get the total number of keys
  INFO                      show store statistics (keys, hits, misses, etc)
  FLUSH                     remove all keys from the store
  HELP                      show this message
  EXIT|QUIT                 leave the REPL`,
	)
}

func printEngineErr(err error) {
	switch {
	case errors.Is(err, engine.ErrEmptyKey):
		fmt.Println("key must not be empty")
	case errors.Is(err, engine.ErrNotFound):
		fmt.Println("key not found")
	case errors.Is(err, engine.ErrValueTooBig):
		fmt.Println("value too large, try shorter value")
	default:
		fmt.Println("error:", err)
	}
}
