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

func TestPeopleIndexShowsAlbumAndBookCounts(t *testing.T) {
	store := &fakeStore{people: []domain.Person{
		{ID: "person-1", Name: "Artist Author", Roles: []string{domain.PersonRoleArtist, domain.PersonRoleAuthor}, AlbumCount: 12, BookCount: 3},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/people", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, `Albumy <strong>12</strong>`) || !strings.Contains(body, `Książki <strong>3</strong>`) {
		t.Fatalf("people index should show album and book counts, body: %s", body)
	}
}

func TestPeopleIndexShowsCountsOnlyForAssignedTypes(t *testing.T) {
	store := &fakeStore{people: []domain.Person{
		{ID: "artist", Name: "Artist", Roles: []string{domain.PersonRoleArtist}, AlbumCount: 4, BookCount: 9},
		{ID: "author", Name: "Author", Roles: []string{domain.PersonRoleAuthor}, AlbumCount: 8, BookCount: 3},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/people", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if strings.Contains(body, `Książki <strong>9</strong>`) || strings.Contains(body, `Albumy <strong>8</strong>`) {
		t.Fatalf("people index should hide counts for unassigned types, body: %s", body)
	}
	if !strings.Contains(body, `data-filter-types="artist `) || !strings.Contains(body, `data-filter-types="author `) {
		t.Fatalf("people index should expose types for filtering, body: %s", body)
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
	if !strings.Contains(body, `aria-label="Edytuj"`) {
		t.Fatalf("catalog edit links should have accessible labels, body: %s", body)
	}
	if strings.Contains(body, `type="submit">Filtruj`) {
		t.Fatalf("catalog should not show a redundant filter submit button, body: %s", body)
	}
}

func TestCatalogShowsAddMenu(t *testing.T) {
	server := NewServer(config.Config{}, &fakeStore{})
	request := authenticatedRequest(server, http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	for _, path := range []string{"/music/new", "/books/new", "/people/new", "/groups/new"} {
		if !strings.Contains(body, `href="`+path+`"`) {
			t.Errorf("add menu should link to %q", path)
		}
	}
	if !strings.Contains(body, `<span>Dodaj</span>`) {
		t.Fatalf("catalog should show a single add menu trigger, body: %s", body)
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

func TestCatalogNormalizesLegacySmallGridView(t *testing.T) {
	store := &fakeStore{}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/items?view=grid-small", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if store.lastFilter.View != "grid" {
		t.Fatalf("view = %q, want grid", store.lastFilter.View)
	}
}

func TestCatalogShowsIconViewSwitch(t *testing.T) {
	server := NewServer(config.Config{}, &fakeStore{})
	request := authenticatedRequest(server, http.MethodGet, "/items?view=list", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if strings.Contains(body, `<select name="view">`) {
		t.Fatalf("catalog should not show a view select, body: %s", body)
	}
	if !strings.Contains(body, `type="radio" name="view" value="grid"`) {
		t.Fatalf("catalog should show a grid view option, body: %s", body)
	}
	if !strings.Contains(body, `type="radio" name="view" value="list" checked`) {
		t.Fatalf("catalog should mark the list view option as checked, body: %s", body)
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != catalogViewCookie || cookies[0].Value != "list" {
		t.Fatalf("unexpected catalog view cookies: %#v", cookies)
	}
}

func TestCatalogUsesRememberedView(t *testing.T) {
	store := &fakeStore{}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: catalogViewCookie, Value: "list"})
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if store.lastFilter.View != "list" {
		t.Fatalf("view = %q, want list", store.lastFilter.View)
	}
	if !strings.Contains(response.Body.String(), `type="radio" name="view" value="list" checked`) {
		t.Fatalf("catalog should mark the remembered list view as checked, body: %s", response.Body.String())
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

func TestGroupFormShowsOnlyArtistsAsMembers(t *testing.T) {
	store := &fakeStore{people: []domain.Person{
		{ID: "artist", Name: "Artist Person", Roles: []string{domain.PersonRoleArtist}},
		{ID: "author", Name: "Author Person", Roles: []string{domain.PersonRoleAuthor}},
		{ID: "both", Name: "Both Person", Roles: []string{domain.PersonRoleArtist, domain.PersonRoleAuthor}},
	}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/groups/new", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, `value="artist"`) || !strings.Contains(body, `value="both"`) {
		t.Fatalf("group form should show artists, body: %s", body)
	}
	if strings.Contains(body, `value="author"`) {
		t.Fatalf("group form should not show author-only people, body: %s", body)
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
