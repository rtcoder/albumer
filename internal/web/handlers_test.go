package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"albumer/internal/config"
	"albumer/internal/domain"
)

type fakeStore struct {
	items      []domain.CatalogItem
	musicRefs  []domain.MusicReference
	people     []domain.Person
	groups     []domain.Group
	lastFilter domain.CatalogFilter
}

func (f *fakeStore) EnsureSchema(ctx context.Context) error { return nil }
func (f *fakeStore) Close(ctx context.Context) error        { return nil }

func (f *fakeStore) ListCatalog(ctx context.Context, filter domain.CatalogFilter) ([]domain.CatalogItem, error) {
	f.lastFilter = filter
	filtered := make([]domain.CatalogItem, 0)
	for _, item := range f.items {
		if filter.Medium != "" && item.Medium != filter.Medium {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered, nil
}

func (f *fakeStore) ListPeople(ctx context.Context) ([]domain.Person, error) { return f.people, nil }
func (f *fakeStore) ListGroups(ctx context.Context) ([]domain.Group, error)  { return f.groups, nil }
func (f *fakeStore) ListMusicReferences(ctx context.Context, query, excludeID string) ([]domain.MusicReference, error) {
	var items []domain.MusicReference
	for _, item := range f.musicRefs {
		if item.ID != excludeID {
			items = append(items, item)
		}
	}
	return items, nil
}
func (f *fakeStore) GetPerson(ctx context.Context, id string) (domain.Person, error) {
	return domain.Person{ID: id, Name: "Person"}, nil
}
func (f *fakeStore) GetGroup(ctx context.Context, id string) (domain.Group, error) {
	return domain.Group{ID: id, Name: "Group"}, nil
}
func (f *fakeStore) GetMusicItem(ctx context.Context, id string) (domain.MusicItem, error) {
	return domain.MusicItem{ID: id, Title: "Music", Medium: domain.MediumCD, Artists: []domain.Contributor{{ID: "p1", Kind: "person"}}}, nil
}
func (f *fakeStore) GetBook(ctx context.Context, id string) (domain.Book, error) {
	return domain.Book{ID: id, Title: "Book", Authors: []domain.Person{{ID: "p1"}}}, nil
}
func (f *fakeStore) CreatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	return person, nil
}
func (f *fakeStore) CreateGroup(ctx context.Context, group domain.Group, memberIDs []string) (domain.Group, error) {
	return group, nil
}
func (f *fakeStore) CreateMusicItem(ctx context.Context, item domain.MusicItem) (domain.MusicItem, error) {
	return item, nil
}
func (f *fakeStore) CreateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	return book, nil
}
func (f *fakeStore) UpdatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	return person, nil
}
func (f *fakeStore) UpdateGroup(ctx context.Context, group domain.Group, memberIDs []string) (domain.Group, error) {
	return group, nil
}
func (f *fakeStore) UpdateMusicItem(ctx context.Context, item domain.MusicItem) (domain.MusicItem, error) {
	return item, nil
}
func (f *fakeStore) UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	return book, nil
}

func TestCatalogRequiresLogin(t *testing.T) {
	server := NewServer(config.Config{}, &fakeStore{})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusSeeOther)
	}
	if location := response.Header().Get("Location"); location != "/login" {
		t.Fatalf("Location = %q, want /login", location)
	}
}

func TestCatalogShowsMusicMedium(t *testing.T) {
	store := &fakeStore{items: []domain.CatalogItem{
		{ID: "1", Kind: "music", Title: "Same album", Medium: domain.MediumCD},
		{ID: "2", Kind: "music", Title: "Same album", Medium: domain.MediumVinyl},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, "CD") || !strings.Contains(body, "VINYL") {
		t.Fatalf("catalog should show separate media, body: %s", body)
	}
}

func TestCatalogShowsDefaultCoverImages(t *testing.T) {
	store := &fakeStore{items: []domain.CatalogItem{
		{ID: "book", Kind: "book", Title: "Book"},
		{ID: "cd", Kind: "music", Title: "CD", Medium: domain.MediumCD},
		{ID: "vinyl", Kind: "music", Title: "Vinyl", Medium: domain.MediumVinyl},
		{ID: "cassette", Kind: "music", Title: "Cassette", Medium: domain.MediumCassette},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	for _, path := range []string{
		"/assets/default-book.svg",
		"/assets/default-cd.svg",
		"/assets/default-vinyl.svg",
		"/assets/default-cassette.svg",
	} {
		if !strings.Contains(body, path) {
			t.Errorf("catalog should show default cover %q", path)
		}
	}
}

func TestCatalogShowsEditLinks(t *testing.T) {
	store := &fakeStore{items: []domain.CatalogItem{
		{ID: "hqw(=.-!d#", Kind: "music", Title: "Album", Medium: domain.MediumCD},
		{ID: "%TZ.a7XrPd", Kind: "book", Title: "Book"},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, `/music/hqw%28=.-%21d%23/edit`) || !strings.Contains(body, `/books/%25TZ.a7XrPd/edit`) {
		t.Fatalf("catalog should show edit links, body: %s", body)
	}
}

func TestCatalogPassesMediumFilter(t *testing.T) {
	store := &fakeStore{items: []domain.CatalogItem{
		{ID: "1", Kind: "music", Title: "CD copy", Medium: domain.MediumCD},
		{ID: "2", Kind: "music", Title: "Vinyl copy", Medium: domain.MediumVinyl},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/items?medium=VINYL", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if store.lastFilter.Medium != domain.MediumVinyl {
		t.Fatalf("medium filter = %q, want VINYL", store.lastFilter.Medium)
	}
	body := response.Body.String()
	if strings.Contains(body, "CD copy") || !strings.Contains(body, "Vinyl copy") {
		t.Fatalf("unexpected filtered body: %s", body)
	}
}

func TestMusicSuggestionsShowOtherMedia(t *testing.T) {
	store := &fakeStore{musicRefs: []domain.MusicReference{
		{ID: "cd-1", Title: "Album", Medium: domain.MediumCD},
		{ID: "cassette-1", Title: "Album", Medium: domain.MediumCassette},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/music/suggestions?q=album&exclude=cd-1", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if strings.Contains(body, "cd-1") || !strings.Contains(body, "cassette-1") || !strings.Contains(body, "CASSETTE") {
		t.Fatalf("unexpected suggestions body: %s", body)
	}
}

func TestPeopleIndexShowsSearch(t *testing.T) {
	store := &fakeStore{people: []domain.Person{{ID: "p1", Name: "O.S.T.R", SortName: "OSTR"}}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/people", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, `aria-label="Szukaj osób"`) || !strings.Contains(body, `data-filter-label="O.S.T.R OSTR"`) {
		t.Fatalf("people index should show searchable rows, body: %s", body)
	}
}

func TestGroupsIndexShowsSearch(t *testing.T) {
	store := &fakeStore{groups: []domain.Group{{ID: "g1", Name: "Płomień 81", SortName: "Plomien 81"}}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/groups", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, `aria-label="Szukaj zespołów"`) || !strings.Contains(body, `data-filter-label="Płomień 81 Plomien 81"`) {
		t.Fatalf("groups index should show searchable rows, body: %s", body)
	}
}

func authenticatedRequest(server *Server, method, target string, form url.Values) *http.Request {
	var body *strings.Reader
	if form == nil {
		body = strings.NewReader("")
	} else {
		body = strings.NewReader(form.Encode())
	}
	request := httptest.NewRequest(method, target, body)
	if form != nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	sessionID := newToken()
	server.sessions[sessionID] = session{CSRF: "csrf-test"}
	request.AddCookie(&http.Cookie{Name: sessionCookie, Value: sessionID})
	return request
}
