package importer

import "testing"

func TestLegacyStatus(t *testing.T) {
	tests := []struct {
		value int
		want  string
	}{
		{1, "owned"},
		{0, ""},
		{7, "legacy:7"},
	}

	for _, tt := range tests {
		if got := LegacyStatus(tt.value); got != tt.want {
			t.Fatalf("LegacyStatus(%d) = %q, want %q", tt.value, got, tt.want)
		}
	}
}
