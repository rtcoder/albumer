package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"albumer/internal/config"
	"albumer/internal/domain"
)

func TestLanguageDefaultsToPolish(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	if language := languageFromRequest(request); language != "pl" {
		t.Fatalf("language = %q, want pl", language)
	}
}

func TestLanguageUsesCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: languageCookie, Value: "en"})

	if language := languageFromRequest(request); language != "en" {
		t.Fatalf("language = %q, want en", language)
	}
}

func TestTranslationFilesHaveMatchingKeys(t *testing.T) {
	for key := range messages["pl"] {
		if _, ok := messages["en"][key]; !ok {
			t.Errorf("English translations are missing key %q", key)
		}
	}
	for key := range messages["en"] {
		if _, ok := messages["pl"][key]; !ok {
			t.Errorf("Polish translations are missing key %q", key)
		}
	}
}

func TestSetLanguageStoresCookieAndRedirects(t *testing.T) {
	server := NewServer(config.Config{}, &fakeStore{})
	request := httptest.NewRequest(http.MethodGet, "/language/en?next=/people", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusSeeOther)
	}
	if location := response.Header().Get("Location"); location != "/people" {
		t.Fatalf("Location = %q, want /people", location)
	}
	if cookie := response.Result().Cookies()[0]; cookie.Name != languageCookie || cookie.Value != "en" {
		t.Fatalf("unexpected language cookie: %#v", cookie)
	}
}

func TestSetLanguageRejectsExternalRedirect(t *testing.T) {
	server := NewServer(config.Config{}, &fakeStore{})
	request := httptest.NewRequest(http.MethodGet, "/language/en?next=//example.com", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if location := response.Header().Get("Location"); location != "/" {
		t.Fatalf("Location = %q, want /", location)
	}
}

func TestPeopleIndexCanRenderEnglish(t *testing.T) {
	store := &fakeStore{people: []domain.Person{{ID: "p1", Name: "O.S.T.R", SortName: "OSTR"}}}
	server := NewServer(config.Config{}, store)
	request := authenticatedRequest(server, http.MethodGet, "/people", nil)
	request.AddCookie(&http.Cookie{Name: languageCookie, Value: "en"})
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, `lang="en"`) || !strings.Contains(body, `aria-label="Search people"`) {
		t.Fatalf("people index should render in English, body: %s", body)
	}
}
