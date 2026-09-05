package cli

import (
	"fmt"
	"time"
)

func (r *REPL) cmdSet(args []string) {
	if len(args) < 2 {
		fmt.Println("usage: SET <key> <value> [ttl] (30s, 5m, 1h)")
		return
	}

	key, value := args[0], args[1]
	var ttl time.Duration

	if len(args) >= 3 {
		parsed, err := time.ParseDuration(args[2])
		if err != nil {
			fmt.Printf("invalid ttl %s: %v\n", args[2], err)
			return
		}
		ttl = parsed
	}

	if err := r.engine.Set(key, []byte(value), ttl); err != nil {
		printEngineErr(err)
		return
	}

	fmt.Println("OK")
}

func (r *REPL) cmdGet(args []string) {
	if len(args) != 1 {
		fmt.Println("usage: GET <key>")
		return
	}

	value, err := r.engine.Get(args[0])
	if err != nil {
		printEngineErr(err)
		return
	}
	fmt.Println(string(value))
}

func (r *REPL) cmdDelete(args []string) {
	if len(args) != 1 {
		fmt.Println("usage: DELETE <key>")
		return
	}

	if err := r.engine.Delete(args[0]); err != nil {
		printEngineErr(err)
		return
	}
	fmt.Println("OK")
}

func (r *REPL) cmdExists(args []string) {
	if len(args) != 1 {
		fmt.Println("usage: EXISTS <key>")
		return
	}

	exists, err := r.engine.Exists(args[0])
	if err != nil {
		printEngineErr(err)
		return
	}
	fmt.Println(exists)
}

func (r *REPL) cmdCount() {
	count := r.engine.Count()
	fmt.Println(count)
}

func (r *REPL) cmdInfo() {
	info := r.engine.Info()

	fmt.Printf("%-10s %d\n", "keys:", info.KeyCount)
	fmt.Printf("%-10s %d\n", "sets:", info.Sets)
	fmt.Printf("%-10s %d\n", "deletes:", info.Deletes)
	fmt.Printf("%-10s %d\n", "hits:", info.Hits)
	fmt.Printf("%-10s %d\n", "misses:", info.Misses)
	fmt.Printf("%-10s %d\n", "evictions:", info.Evictions)

	total := info.Hits + info.Misses
	if total > 0 {
		fmt.Printf(
			"%-10s %.1f%%\n",
			"hit rate:",
			float64(info.Hits)/float64(total)*100,
		)
	}
}
