package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"kivo/engine"
	"kivo/internal/server/config"
	"kivo/internal/server/handler"
	"kivo/internal/server/middleware"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGABRT,
	)
	defer stop()

	eng, err := engine.New(ctx, engine.DefaultOpts())
	if err != nil {
		log.Fatalf("failed to start the engine: %v\n", err)
	}

	cfg := config.Load()

	mux := http.NewServeMux()

	hdl := handler.New(eng)
	mux.HandleFunc("GET /ping", hdl.Ping)

	mux.HandleFunc("PUT /kv/{key}", hdl.Put)
	mux.HandleFunc("DELETE /kv/{key}", hdl.Delete)

	mux.HandleFunc("GET /kv/{key}", hdl.Get)
	mux.HandleFunc("GET /kv/{key}/exists", hdl.Exists)
	mux.HandleFunc("GET /stats", hdl.Info)

	mux.HandleFunc("GET /key/{key}/ttl", hdl.TTL)
	mux.HandleFunc("PUT /key/{key}/expire", hdl.Expire)
	mux.HandleFunc("PUT /key/{key}/persist", hdl.Persist)

	srv := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: middleware.Logger(mux),
	}

	go func() {
		log.Println("server starting on port", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("listen error: %v", err)
			}
		}
	}()

	<-ctx.Done()

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
