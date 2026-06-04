package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const themeCookie = "albumer_theme"

func normalizeTheme(theme string) string {
	if strings.EqualFold(theme, "light") {
		return "light"
	}
	return "dark"
}

func themeFromRequest(r *http.Request) string {
	if cookie, err := r.Cookie(themeCookie); err == nil {
		return normalizeTheme(cookie.Value)
	}
	return "dark"
}

func (s *Server) setTheme(w http.ResponseWriter, r *http.Request) {
	theme := normalizeTheme(chi.URLParam(r, "theme"))
	http.SetCookie(w, &http.Cookie{
		Name:     themeCookie,
		Value:    theme,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})

	next := r.URL.Query().Get("next")
	if !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		next = "/"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}
