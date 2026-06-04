package domain

import "testing"

func TestMusicItemHasRelated(t *testing.T) {
	item := MusicItem{Related: []MusicReference{{ID: "other"}}}
	if !item.HasRelated("other") {
		t.Fatal("expected related item")
	}
	if item.HasRelated("missing") {
		t.Fatal("did not expect missing item")
	}
}
