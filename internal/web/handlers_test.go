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

func (f *fakeStore) ListPeople(ctx context.Context) ([]domain.Person, error) { return nil, nil }
func (f *fakeStore) ListGroups(ctx context.Context) ([]domain.Group, error)  { return nil, nil }
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
