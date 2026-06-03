package web

import (
	"net/http"
	"strconv"
	"strings"

	"albumer/internal/domain"
)

type pageData struct {
	Title  string
	CSRF   string
	Error  string
	Filter domain.CatalogFilter
	Items  []domain.CatalogItem
	People []domain.Person
	Groups []domain.Group
}

func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {
	filter := domain.CatalogFilter{
		Query:  strings.TrimSpace(r.URL.Query().Get("q")),
		Kind:   r.URL.Query().Get("kind"),
		Medium: domain.NormalizeMedium(r.URL.Query().Get("medium")),
		Status: r.URL.Query().Get("status"),
		View:   r.URL.Query().Get("view"),
	}
	if filter.View == "" {
		filter.View = "grid-large"
	}
	if !filter.Medium.Valid() {
		filter.Medium = ""
	}

	items, err := s.store.ListCatalog(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := pageData{Title: "Katalog", CSRF: s.csrf(r), Filter: filter, Items: items}
	if r.Header.Get("HX-Request") == "true" || r.Header.Get("X-Fragment") == "catalog" {
		s.render(w, http.StatusOK, "catalog_items", data)
		return
	}
	s.render(w, http.StatusOK, "catalog", data)
}

func (s *Server) musicForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.formData(r, "Dodaj muzyke")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, http.StatusOK, "music_form", data)
}

func (s *Server) createMusic(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		s.renderFormError(w, r, "music_form", "Nieprawidlowy formularz.")
		return
	}

	coverPath, err := s.coverFromRequest(r)
	if err != nil {
		s.renderFormError(w, r, "music_form", err.Error())
		return
	}

	item := domain.MusicItem{
		Title:     strings.TrimSpace(r.FormValue("title")),
		Medium:    domain.NormalizeMedium(r.FormValue("medium")),
		Year:      parseYear(r.FormValue("year")),
		Status:    strings.TrimSpace(r.FormValue("status")),
		CoverPath: coverPath,
		Notes:     strings.TrimSpace(r.FormValue("notes")),
		Artists:   contributorsFromForm(r),
	}
	if _, err := s.store.CreateMusicItem(r.Context(), item); err != nil {
		s.renderFormError(w, r, "music_form", err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) bookForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.formData(r, "Dodaj ksiazke")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, http.StatusOK, "book_form", data)
}

func (s *Server) createBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		s.renderFormError(w, r, "book_form", "Nieprawidlowy formularz.")
		return
	}

	coverPath, err := s.coverFromRequest(r)
	if err != nil {
		s.renderFormError(w, r, "book_form", err.Error())
		return
	}

	authors := make([]domain.Person, 0)
	for _, id := range r.Form["person_ids"] {
		authors = append(authors, domain.Person{ID: id})
	}

	book := domain.Book{
		Title:     strings.TrimSpace(r.FormValue("title")),
		ISBN:      strings.TrimSpace(r.FormValue("isbn")),
		Year:      parseYear(r.FormValue("year")),
		Status:    strings.TrimSpace(r.FormValue("status")),
		CoverPath: coverPath,
		Notes:     strings.TrimSpace(r.FormValue("notes")),
		Authors:   authors,
	}
	if _, err := s.store.CreateBook(r.Context(), book); err != nil {
		s.renderFormError(w, r, "book_form", err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) personForm(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "person_form", pageData{Title: "Dodaj osobe", CSRF: s.csrf(r)})
}

func (s *Server) createPerson(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.render(w, http.StatusBadRequest, "person_form", pageData{Title: "Dodaj osobe", CSRF: s.csrf(r), Error: "Nieprawidlowy formularz."})
		return
	}
	person := domain.Person{
		Name:     strings.TrimSpace(r.FormValue("name")),
		SortName: strings.TrimSpace(r.FormValue("sort_name")),
		Notes:    strings.TrimSpace(r.FormValue("notes")),
	}
	if _, err := s.store.CreatePerson(r.Context(), person); err != nil {
		s.render(w, http.StatusBadRequest, "person_form", pageData{Title: "Dodaj osobe", CSRF: s.csrf(r), Error: err.Error()})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) groupForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.formData(r, "Dodaj zespol")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, http.StatusOK, "group_form", data)
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderFormError(w, r, "group_form", "Nieprawidlowy formularz.")
		return
	}
	group := domain.Group{
		Name:     strings.TrimSpace(r.FormValue("name")),
		SortName: strings.TrimSpace(r.FormValue("sort_name")),
		Notes:    strings.TrimSpace(r.FormValue("notes")),
	}
	if _, err := s.store.CreateGroup(r.Context(), group, r.Form["member_ids"]); err != nil {
		s.renderFormError(w, r, "group_form", err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) formData(r *http.Request, title string) (pageData, error) {
	people, err := s.store.ListPeople(r.Context())
	if err != nil {
		return pageData{}, err
	}
	groups, err := s.store.ListGroups(r.Context())
	if err != nil {
		return pageData{}, err
	}
	return pageData{Title: title, CSRF: s.csrf(r), People: people, Groups: groups}, nil
}

func (s *Server) renderFormError(w http.ResponseWriter, r *http.Request, templateName, message string) {
	data, err := s.formData(r, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Error = message
	s.render(w, http.StatusBadRequest, templateName, data)
}

func (s *Server) coverFromRequest(r *http.Request) (string, error) {
	file, header, err := r.FormFile("cover")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer file.Close()
	return s.saveCover(file, header)
}

func (s *Server) csrf(r *http.Request) string {
	sess, ok := s.currentSession(r)
	if !ok {
		return ""
	}
	return sess.CSRF
}

func parseYear(value string) int {
	year, _ := strconv.Atoi(strings.TrimSpace(value))
	return year
}

func contributorsFromForm(r *http.Request) []domain.Contributor {
	contributors := make([]domain.Contributor, 0)
	for _, id := range r.Form["person_ids"] {
		contributors = append(contributors, domain.Contributor{ID: id, Kind: "person"})
	}
	for _, id := range r.Form["group_ids"] {
		contributors = append(contributors, domain.Contributor{ID: id, Kind: "group"})
	}
	return contributors
}
