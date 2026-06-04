package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"albumer/internal/domain"

	"github.com/go-chi/chi/v5"
)

const languageCookie = "albumer_language"

var messages = loadMessages("pl", "en")

var errorMessageKeys = map[string]string{
	"title is required":                        "error.title_required",
	"medium must be CD, VINYL, or CASSETTE":    "error.medium_invalid",
	"status must be owned, ordered, or wanted": "error.status_invalid",
	"at least one artist or group is required": "error.contributor_required",
	"at least one author is required":          "error.author_required",
	"name is required":                         "error.name_required",
	"at least one person type is required":     "error.person_type_required",
	"person type must be artist or author":     "error.person_type_invalid",
	"cover must be JPEG, PNG, or WebP":         "error.cover_invalid",
}

func loadMessages(languages ...string) map[string]map[string]string {
	loaded := make(map[string]map[string]string, len(languages))
	for _, language := range languages {
		path := "locales/" + language + ".json"
		content, err := webFiles.ReadFile(path)
		if err != nil {
			panic(fmt.Errorf("read translations %s: %w", path, err))
		}

		var translations map[string]string
		if err := json.Unmarshal(content, &translations); err != nil {
			panic(fmt.Errorf("parse translations %s: %w", path, err))
		}
		loaded[language] = translations
	}
	return loaded
}

func normalizeLanguage(language string) string {
	if strings.EqualFold(language, "en") || strings.HasPrefix(strings.ToLower(language), "en-") {
		return "en"
	}
	return "pl"
}

func languageFromRequest(r *http.Request) string {
	if language := r.URL.Query().Get("lang"); language != "" {
		return normalizeLanguage(language)
	}
	if cookie, err := r.Cookie(languageCookie); err == nil {
		return normalizeLanguage(cookie.Value)
	}
	return "pl"
}

func translate(language, key string) string {
	language = normalizeLanguage(language)
	if value := messages[language][key]; value != "" {
		return value
	}
	if value := messages["pl"][key]; value != "" {
		return value
	}
	return key
}

func translateError(language, message string) string {
	if key := errorMessageKeys[message]; key != "" {
		return translate(language, key)
	}
	if strings.HasPrefix(message, "error.") {
		return translate(language, message)
	}
	return message
}

func (d pageData) Tr(key string) string {
	return translate(d.Lang, key)
}

func (d pageData) StatusLabel(status domain.ItemStatus) string {
	return d.Tr("status." + string(status))
}

func (s *Server) setLanguage(w http.ResponseWriter, r *http.Request) {
	language := normalizeLanguage(chi.URLParam(r, "lang"))
	http.SetCookie(w, &http.Cookie{
		Name:     languageCookie,
		Value:    language,
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
