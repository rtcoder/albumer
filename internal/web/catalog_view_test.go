package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCatalogViewUsesCookieWithoutQueryParameter(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: catalogViewCookie, Value: "list"})

	view, selected := catalogViewFromRequest(request)

	if view != "list" {
		t.Fatalf("view = %q, want list", view)
	}
	if selected {
		t.Fatal("cookie view should not be treated as a newly selected view")
	}
}

func TestCatalogViewQueryOverridesCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/items?view=grid", nil)
	request.AddCookie(&http.Cookie{Name: catalogViewCookie, Value: "list"})

	view, selected := catalogViewFromRequest(request)

	if view != "grid" {
		t.Fatalf("view = %q, want grid", view)
	}
	if !selected {
		t.Fatal("query view should be treated as a newly selected view")
	}
}
