package main

import (
	"context"
	"flag"
	"log"

	"albumer/internal/config"
	"albumer/internal/importer"
	"albumer/internal/storage"
)

func main() {
	path := flag.String("file", "json/data.json", "path to legacy data.json")
	flag.Parse()

	cfg := config.Load()

	ctx := context.Background()
	store, err := storage.NewNeo4jStore(ctx, cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword)
	if err != nil {
		log.Fatalf("connect neo4j: %v", err)
	}
	defer store.Close(ctx)

	if err := store.EnsureSchema(ctx); err != nil {
		log.Fatalf("ensure neo4j schema: %v", err)
	}

	data, err := importer.LoadLegacyData(*path)
	if err != nil {
		log.Fatalf("load legacy json: %v", err)
	}

	summary, err := importer.ImportLegacyCDs(ctx, store.Driver(), data)
	if err != nil {
		log.Fatalf("import legacy cds: %v", err)
	}

	log.Printf("imported people=%d groups=%d group_links=%d cds=%d cd_links=%d",
		summary.People, summary.Groups, summary.GroupLinks, summary.CDs, summary.CDLinks)
}
