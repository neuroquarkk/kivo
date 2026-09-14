package config

import "os"

type Env struct {
	PORT string
}

func Load() *Env {
	return &Env{
		PORT: lookup("PORT", "8081"),
	}
}

func lookup(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
