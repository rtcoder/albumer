package config

import (
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Addr              string
	Neo4jURI          string
	Neo4jUser         string
	Neo4jPassword     string
	AppSecret         string
	UploadDir         string
	AdminPasswordHash string
}

func Load() Config {
	dotenv, err := ReadDotenv(FindDotenv("."))
	if err != nil {
		log.Fatalf("load .env: %v", err)
	}

	return Config{
		Addr:              env(dotenv, "ADDR", ":4000"),
		Neo4jURI:          env(dotenv, "NEO4J_URI", "neo4j://localhost:7687"),
		Neo4jUser:         env(dotenv, "NEO4J_USER", "neo4j"),
		Neo4jPassword:     env(dotenv, "NEO4J_PASSWORD", "neo4j"),
		AppSecret:         env(dotenv, "APP_SECRET", "dev-secret-change-me"),
		UploadDir:         env(dotenv, "UPLOAD_DIR", filepath.Join(".", "uploads")),
		AdminPasswordHash: env(dotenv, "ADMIN_PASSWORD_HASH", ""),
	}
}

func env(dotenv map[string]string, key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = dotenv[key]
	}
	if value == "" {
		value = fallback
	}
	return value
}
