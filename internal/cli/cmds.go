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

	printStat := func(name string, value any) {
		fmt.Printf("%-20s %v\n", name+":", value)
	}

	printStat("keys", info.KeyCount)
	printStat("sets", info.Sets)
	printStat("deletes", info.Deletes)
	printStat("hits", info.Hits)
	printStat("misses", info.Misses)
	printStat("evictions", info.Evictions)
	printStat("memory", info.MemLimit)
	printStat("memory per bucket", info.MemPerBucket)

	total := info.Hits + info.Misses
	if total > 0 {
		printStat("hit rate", fmt.Sprintf("%.1f%%",
			float64(info.Hits)/float64(total)*100,
		))
	}
}

func (r *REPL) cmdFlush() {
	r.engine.Flush()
	fmt.Println("flushed all keys")
}

func (r *REPL) cmdTTL(args []string) {
	if len(args) != 1 {
		fmt.Println("usage: TTL <key>")
		return
	}

	remaining, hasTTL, err := r.engine.TTL(args[0])
	if err != nil {
		printEngineErr(err)
		return
	}

	if !hasTTL {
		fmt.Println("permanent")
		return
	}

	fmt.Println(remaining.Round(time.Second).String())
}

func (r *REPL) cmdExpire(args []string) {
	if len(args) != 2 {
		fmt.Println("usage: EXPIRE <key> <ttl>")
		return
	}

	parsed, err := time.ParseDuration(args[1])
	if err != nil {
		fmt.Printf("invalid ttl %s: %v\n", args[1], err)
		return
	}

	if err := r.engine.Expire(args[0], parsed); err != nil {
		printEngineErr(err)
		return
	}
	fmt.Println("OK")
}

func (r *REPL) cmdPersist(args []string) {
	if len(args) != 1 {
		fmt.Println("usage: PERSIST <key>")
		return
	}

	if err := r.engine.Persist(args[0]); err != nil {
		printEngineErr(err)
		return
	}
	fmt.Println("OK")
}
