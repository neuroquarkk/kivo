package main

import (
	"context"
	"log"
	"net/http"

	"kivo/engine"
	"kivo/internal/server/config"
	"kivo/internal/server/handler"
	"kivo/internal/server/middleware"
)

func main() {
	ctx := context.Background()
	eng, err := engine.New(ctx, engine.DefaultOpts())
	if err != nil {
		log.Fatalf("failed to start the engine: %v\n", err)
	}

	cfg := config.Load()

	mux := http.NewServeMux()

	hdl := handler.New(eng)
	mux.HandleFunc("GET /ping", hdl.Ping)
	mux.HandleFunc("PUT /kv/{key}", hdl.Put)

	srv := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: middleware.Logger(mux),
	}

	log.Println("server starting on port", srv.Addr)
	srv.ListenAndServe()
}
