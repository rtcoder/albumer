package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"albumer/internal/config"
)

func TestThemeDefaultsToDark(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	if theme := themeFromRequest(request); theme != "dark" {
		t.Fatalf("theme = %q, want dark", theme)
	}
}

func TestThemeUsesCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: themeCookie, Value: "light"})

	if theme := themeFromRequest(request); theme != "light" {
		t.Fatalf("theme = %q, want light", theme)
	}
}

func TestSetThemeStoresCookieAndRedirects(t *testing.T) {
	server := NewServer(config.Config{}, &fakeStore{})
	request := httptest.NewRequest(http.MethodGet, "/theme/light?next=/people", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusSeeOther)
	}
	if location := response.Header().Get("Location"); location != "/people" {
		t.Fatalf("Location = %q, want /people", location)
	}
	if cookie := response.Result().Cookies()[0]; cookie.Name != themeCookie || cookie.Value != "light" {
		t.Fatalf("unexpected theme cookie: %#v", cookie)
	}
}

func TestLoginCanRenderLightTheme(t *testing.T) {
	server := NewServer(config.Config{}, &fakeStore{})
	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	request.AddCookie(&http.Cookie{Name: themeCookie, Value: "light"})
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, `data-theme="light"`) || !strings.Contains(body, `aria-label="Jasny motyw"`) {
		t.Fatalf("login should render light theme, body: %s", body)
	}
}
