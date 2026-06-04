package web

import (
	"net/http"
	"time"
)

const catalogViewCookie = "albumer_catalog_view"

func normalizeCatalogView(view string) string {
	if view == "list" {
		return "list"
	}
	return "grid"
}

func catalogViewFromRequest(r *http.Request) (string, bool) {
	if _, exists := r.URL.Query()["view"]; exists {
		return normalizeCatalogView(r.URL.Query().Get("view")), true
	}
	if cookie, err := r.Cookie(catalogViewCookie); err == nil {
		return normalizeCatalogView(cookie.Value), false
	}
	return "grid", false
}

func setCatalogViewCookie(w http.ResponseWriter, view string) {
	http.SetCookie(w, &http.Cookie{
		Name:     catalogViewCookie,
		Value:    normalizeCatalogView(view),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
}
