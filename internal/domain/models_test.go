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
		Status:  StatusOwned,
		Artists: []Contributor{{ID: "p1", Kind: "person", Name: "Artist"}},
	}

	if err := ValidateMusicItem(item); err == nil {
		t.Fatal("expected missing medium error")
	}
}

func TestValidateMusicItemAcceptsSeparatePhysicalMedia(t *testing.T) {
	artists := []Contributor{{ID: "g1", Kind: "group", Name: "Band"}}
	for _, medium := range ValidMusicMedia {
		item := MusicItem{Title: "Same album", Medium: medium, Status: StatusOwned, Artists: artists}
		if err := ValidateMusicItem(item); err != nil {
			t.Fatalf("ValidateMusicItem(%s) error = %v", medium, err)
		}
	}
}

func TestItemStatusValidationAndLabels(t *testing.T) {
	tests := []struct {
		status ItemStatus
		label  string
	}{
		{StatusOwned, "Posiadane"},
		{StatusOrdered, "Zamówione"},
		{StatusWanted, "Chcę kupić"},
	}

	for _, tt := range tests {
		if !tt.status.Valid() {
			t.Fatalf("%q should be valid", tt.status)
		}
		if got := tt.status.Label(); got != tt.label {
			t.Fatalf("%q label = %q, want %q", tt.status, got, tt.label)
		}
	}

	if ItemStatus("other").Valid() {
		t.Fatal("unknown status should be invalid")
	}
}

func TestValidateBookRequiresValidStatus(t *testing.T) {
	book := Book{Title: "Book", Authors: []Person{{ID: "p1"}}}
	if err := ValidateBook(book); err == nil {
		t.Fatal("expected missing status error")
	}
}

func TestValidatePersonRequiresAtLeastOneType(t *testing.T) {
	if err := ValidatePerson(Person{Name: "Person"}); err == nil {
		t.Fatal("expected missing person type error")
	}
	if err := ValidatePerson(Person{Name: "Person", Roles: []string{PersonRoleArtist, PersonRoleAuthor}}); err != nil {
		t.Fatalf("ValidatePerson() error = %v", err)
	}
}
