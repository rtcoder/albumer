package domain

import "testing"

func TestMusicMediumValidation(t *testing.T) {
	tests := []struct {
		name   string
		medium MusicMedium
		valid  bool
	}{
		{"cd", MediumCD, true},
		{"vinyl", MediumVinyl, true},
		{"cassette", MediumCassette, true},
		{"empty", "", false},
		{"unknown", "TAPE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.medium.Valid(); got != tt.valid {
				t.Fatalf("Valid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestValidateMusicItemRequiresMedium(t *testing.T) {
	item := MusicItem{
		Title:   "A title",
		Artists: []Contributor{{ID: "p1", Kind: "person", Name: "Artist"}},
	}

	if err := ValidateMusicItem(item); err == nil {
		t.Fatal("expected missing medium error")
	}
}

func TestValidateMusicItemAcceptsSeparatePhysicalMedia(t *testing.T) {
	artists := []Contributor{{ID: "g1", Kind: "group", Name: "Band"}}
	for _, medium := range ValidMusicMedia {
		item := MusicItem{Title: "Same album", Medium: medium, Artists: artists}
		if err := ValidateMusicItem(item); err != nil {
			t.Fatalf("ValidateMusicItem(%s) error = %v", medium, err)
		}
	}
}
