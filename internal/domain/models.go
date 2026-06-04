package domain

import (
	"errors"
	"strings"
	"time"
)

type MusicMedium string

const (
	MediumCD       MusicMedium = "CD"
	MediumVinyl    MusicMedium = "VINYL"
	MediumCassette MusicMedium = "CASSETTE"
)

var ValidMusicMedia = []MusicMedium{MediumCD, MediumVinyl, MediumCassette}

type Person struct {
	ID        string
	Name      string
	SortName  string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Group struct {
	ID        string
	Name      string
	SortName  string
	Notes     string
	Members   []Person
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MusicItem struct {
	ID        string
	Title     string
	Medium    MusicMedium
	Year      int
	Status    string
	CoverPath string
	Notes     string
	Artists   []Contributor
	Related   []MusicReference
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MusicReference struct {
	ID     string
	Title  string
	Medium MusicMedium
}

type Book struct {
	ID        string
	Title     string
	ISBN      string
	Year      int
	Status    string
	CoverPath string
	Notes     string
	Authors   []Person
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Contributor struct {
	ID   string
	Name string
	Kind string
}

type CatalogItem struct {
	ID           string
	Kind         string
	Title        string
	Subtitle     string
	Medium       MusicMedium
	Status       string
	CoverPath    string
	Year         int
	Contributors []Contributor
}

func (m MusicItem) HasContributor(kind, id string) bool {
	for _, contributor := range m.Artists {
		if contributor.Kind == kind && contributor.ID == id {
			return true
		}
	}
	return false
}

func (m MusicItem) HasRelated(id string) bool {
	for _, related := range m.Related {
		if related.ID == id {
			return true
		}
	}
	return false
}

func (b Book) HasAuthor(id string) bool {
	for _, author := range b.Authors {
		if author.ID == id {
			return true
		}
	}
	return false
}

func (g Group) HasMember(id string) bool {
	for _, member := range g.Members {
		if member.ID == id {
			return true
		}
	}
	return false
}

type CatalogFilter struct {
	Query  string
	Kind   string
	Medium MusicMedium
	Status string
	View   string
}

func (m MusicMedium) Valid() bool {
	switch m {
	case MediumCD, MediumVinyl, MediumCassette:
		return true
	default:
		return false
	}
}

func NormalizeMedium(value string) MusicMedium {
	return MusicMedium(strings.ToUpper(strings.TrimSpace(value)))
}

func ValidateMusicItem(item MusicItem) error {
	if strings.TrimSpace(item.Title) == "" {
		return errors.New("title is required")
	}
	if !item.Medium.Valid() {
		return errors.New("medium must be CD, VINYL, or CASSETTE")
	}
	if len(item.Artists) == 0 {
		return errors.New("at least one artist or group is required")
	}
	return nil
}

func ValidateBook(book Book) error {
	if strings.TrimSpace(book.Title) == "" {
		return errors.New("title is required")
	}
	if len(book.Authors) == 0 {
		return errors.New("at least one author is required")
	}
	return nil
}

func ValidateNamed(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	return nil
}
