package config

import (
	"os"
	"path/filepath"
	"testing"
)

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
