package domain

import "testing"

func TestDefaultSortNameUsesNameWhenSortNameIsEmpty(t *testing.T) {
	if got := DefaultSortName("  O.S.T.R  ", " "); got != "O.S.T.R" {
		t.Fatalf("DefaultSortName() = %q, want %q", got, "O.S.T.R")
	}
}

func TestDefaultSortNameKeepsProvidedValue(t *testing.T) {
	if got := DefaultSortName("O.S.T.R", "OSTR"); got != "OSTR" {
		t.Fatalf("DefaultSortName() = %q, want %q", got, "OSTR")
	}
}
