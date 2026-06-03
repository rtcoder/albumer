package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"albumer/internal/config"
	"albumer/internal/storage"
	"albumer/internal/web"
)

func main() {
	if err := config.LoadDotenv(".env"); err != nil {
		log.Fatalf("load .env: %v", err)
	}

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store, err := storage.NewNeo4jStore(ctx, cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword)
	if err != nil {
		log.Fatalf("connect neo4j: %v", err)
	}
	defer store.Close(ctx)

	if err := store.EnsureSchema(ctx); err != nil {
		log.Fatalf("ensure neo4j schema: %v", err)
	}

	app := web.NewServer(cfg, store)
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           app.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Albumer listening on %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
