package domain

import "testing"

func TestSearchKey(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"O.S.T.R", "ostr"},
		{"ostr", "ostr"},
		{"O.S.T.R.", "ostr"},
		{"Łona i Webber", "lonaiwebber"},
		{"Zażółć gęślą jaźń", "zazolcgeslajazn"},
		{"  100 dni do matury  ", "100dnidomatury"},
	}

	for _, tt := range tests {
		if got := SearchKey(tt.value); got != tt.want {
			t.Fatalf("SearchKey(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}
