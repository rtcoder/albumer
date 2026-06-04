package main

import (
	"context"
	"log"

	"albumer/internal/config"
	"albumer/internal/importer"
	"albumer/internal/storage"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	store, err := storage.NewNeo4jStore(ctx, cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword)
	if err != nil {
		log.Fatalf("connect neo4j: %v", err)
	}
	defer store.Close(ctx)

	summary, err := importer.BackfillSearchKeys(ctx, store.Driver())
	if err != nil {
		log.Fatalf("backfill search keys: %v", err)
	}

	log.Printf("backfilled people=%d groups=%d music=%d books=%d",
		summary.People, summary.Groups, summary.Music, summary.Books)
}
