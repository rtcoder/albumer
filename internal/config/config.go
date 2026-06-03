package config

import (
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
	return Config{
		Addr:              env("ADDR", ":4000"),
		Neo4jURI:          env("NEO4J_URI", "neo4j://localhost:7687"),
		Neo4jUser:         env("NEO4J_USER", "neo4j"),
		Neo4jPassword:     env("NEO4J_PASSWORD", "neo4j"),
		AppSecret:         env("APP_SECRET", "dev-secret-change-me"),
		UploadDir:         env("UPLOAD_DIR", filepath.Join(".", "uploads")),
		AdminPasswordHash: os.Getenv("ADMIN_PASSWORD_HASH"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
