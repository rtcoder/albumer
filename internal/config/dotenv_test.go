package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesDotenvValues(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	content := "ADDR=:9999\nNEO4J_PASSWORD=from-dotenv\nUPLOAD_DIR=./covers\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()

	if cfg.Addr != ":9999" {
		t.Fatalf("Addr = %q, want :9999", cfg.Addr)
	}
	if cfg.Neo4jPassword != "from-dotenv" {
		t.Fatalf("Neo4jPassword = %q, want from-dotenv", cfg.Neo4jPassword)
	}
	if cfg.UploadDir != "./covers" {
		t.Fatalf("UploadDir = %q, want ./covers", cfg.UploadDir)
	}
}

func TestLoadPrefersEnvironmentOverDotenv(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("ADDR=:9999\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADDR", ":1111")

	cfg := Load()

	if cfg.Addr != ":1111" {
		t.Fatalf("Addr = %q, want :1111", cfg.Addr)
	}
}

func TestFindDotenvWalksUp(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "cmd", "albumer")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("ADDR=:9999\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := FindDotenv(nested); got != filepath.Join(dir, ".env") {
		t.Fatalf("FindDotenv() = %q", got)
	}
}

func TestLoadDotenvSetsMissingValues(t *testing.T) {
	t.Setenv("ALBUMER_TEST_EXISTING", "")
	path := filepath.Join(t.TempDir(), ".env")
	content := "ALBUMER_TEST_VALUE=from-file\n# comment\nALBUMER_TEST_QUOTED=\"quoted value\"\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := LoadDotenv(path); err != nil {
		t.Fatalf("LoadDotenv() error = %v", err)
	}

	if got := os.Getenv("ALBUMER_TEST_VALUE"); got != "from-file" {
		t.Fatalf("ALBUMER_TEST_VALUE = %q", got)
	}
	if got := os.Getenv("ALBUMER_TEST_QUOTED"); got != "quoted value" {
		t.Fatalf("ALBUMER_TEST_QUOTED = %q", got)
	}
}

func TestLoadDotenvDoesNotOverrideExistingValues(t *testing.T) {
	t.Setenv("ALBUMER_TEST_KEEP", "from-env")
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("ALBUMER_TEST_KEEP=from-file\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := LoadDotenv(path); err != nil {
		t.Fatalf("LoadDotenv() error = %v", err)
	}

	if got := os.Getenv("ALBUMER_TEST_KEEP"); got != "from-env" {
		t.Fatalf("ALBUMER_TEST_KEEP = %q", got)
	}
}

func TestLoadDotenvIgnoresMissingFile(t *testing.T) {
	if err := LoadDotenv(filepath.Join(t.TempDir(), ".env")); err != nil {
		t.Fatalf("LoadDotenv() error = %v", err)
	}
}
