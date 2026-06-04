package importer

import (
	"testing"

	"albumer/internal/domain"
)

func TestLegacyStatus(t *testing.T) {
	tests := []struct {
		value int
		want  domain.ItemStatus
	}{
		{1, domain.StatusOwned},
		{0, domain.StatusWanted},
		{7, domain.StatusWanted},
	}

	for _, tt := range tests {
		if got := LegacyStatus(tt.value); got != tt.want {
			t.Fatalf("LegacyStatus(%d) = %q, want %q", tt.value, got, tt.want)
		}
	}
}
