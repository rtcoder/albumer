package web

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"albumer/internal/domain"
	"albumer/internal/storage"
	"github.com/go-chi/chi/v5"
)

type pageData struct {
	Title           string
	TitleKey        string
	Lang            string
	NextURL         string
	CSRF            string
	Error           string
	Filter          domain.CatalogFilter
	Items           []domain.CatalogItem
	People          []domain.Person
	Groups          []domain.Group
	MusicReferences []domain.MusicReference
	Music           domain.MusicItem
	Book            domain.Book
	Person          domain.Person
	Group           domain.Group
	Action          string
	Edit            bool
}

func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {
	filter := domain.CatalogFilter{
		Query:  strings.TrimSpace(r.URL.Query().Get("q")),
		Kind:   r.URL.Query().Get("kind"),
		Medium: domain.NormalizeMedium(r.URL.Query().Get("medium")),
		Status: domain.NormalizeStatus(r.URL.Query().Get("status")),
		View:   r.URL.Query().Get("view"),
	}
	if filter.View == "" {
		filter.View = "grid-large"
	}
	if !filter.Medium.Valid() {
		filter.Medium = ""
	}
	if !filter.Status.Valid() {
		filter.Status = ""
	}

	items, err := s.store.ListCatalog(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := pageData{TitleKey: "title.catalog", CSRF: s.csrf(r), Filter: filter, Items: items}
	if r.Header.Get("HX-Request") == "true" || r.Header.Get("X-Fragment") == "catalog" {
		s.render(w, r, http.StatusOK, "catalog_items", data)
		return
	}
	s.render(w, r, http.StatusOK, "catalog", data)
}

func (s *Server) musicForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.formData(r, "title.add_music", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, http.StatusOK, "music_form", data)
}

func (s *Server) createMusic(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		s.renderFormError(w, r, "music_form", "error.invalid_form")
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
		Status:    domain.NormalizeStatus(r.FormValue("status")),
		CoverPath: coverPath,
		Notes:     strings.TrimSpace(r.FormValue("notes")),
		Artists:   contributorsFromForm(r),
		Related:   relatedFromForm(r),
	}
	if _, err := s.store.CreateMusicItem(r.Context(), item); err != nil {
		s.renderFormError(w, r, "music_form", err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) editMusicForm(w http.ResponseWriter, r *http.Request) {
	id := routeID(r, "id")
	item, err := s.store.GetMusicItem(r.Context(), id)
	if err != nil {
		s.renderNotFound(w, r, err)
		return
	}
	data, err := s.formData(r, "title.edit_music", item.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Music = item
	data.Action = editPath("music", item.ID)
	data.Edit = true
	s.render(w, r, http.StatusOK, "music_form", data)
}

func (s *Server) updateMusic(w http.ResponseWriter, r *http.Request) {
	existing, err := s.store.GetMusicItem(r.Context(), routeID(r, "id"))
	if err != nil {
		s.renderNotFound(w, r, err)
		return
	}
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		s.renderMusicError(w, r, existing, "error.invalid_form")
		return
	}

	coverPath, err := s.coverFromRequest(r)
	if err != nil {
		s.renderMusicError(w, r, existing, err.Error())
		return
	}
	if coverPath == "" {
		coverPath = existing.CoverPath
	}

	item := domain.MusicItem{
		ID:        existing.ID,
		Title:     strings.TrimSpace(r.FormValue("title")),
		Medium:    domain.NormalizeMedium(r.FormValue("medium")),
		Year:      parseYear(r.FormValue("year")),
		Status:    domain.NormalizeStatus(r.FormValue("status")),
		CoverPath: coverPath,
		Notes:     strings.TrimSpace(r.FormValue("notes")),
		Artists:   contributorsFromForm(r),
		Related:   relatedFromForm(r),
	}
	if _, err := s.store.UpdateMusicItem(r.Context(), item); err != nil {
		s.renderMusicError(w, r, item, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) bookForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.formData(r, "title.add_book", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, http.StatusOK, "book_form", data)
}

func (s *Server) createBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		s.renderFormError(w, r, "book_form", "error.invalid_form")
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
		Status:    domain.NormalizeStatus(r.FormValue("status")),
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

func (s *Server) editBookForm(w http.ResponseWriter, r *http.Request) {
	id := routeID(r, "id")
	book, err := s.store.GetBook(r.Context(), id)
	if err != nil {
		s.renderNotFound(w, r, err)
		return
	}
	data, err := s.formData(r, "title.edit_book", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Book = book
	data.Action = editPath("books", book.ID)
	data.Edit = true
	s.render(w, r, http.StatusOK, "book_form", data)
}

func (s *Server) updateBook(w http.ResponseWriter, r *http.Request) {
	existing, err := s.store.GetBook(r.Context(), routeID(r, "id"))
	if err != nil {
		s.renderNotFound(w, r, err)
		return
	}
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		s.renderBookError(w, r, existing, "error.invalid_form")
		return
	}

	coverPath, err := s.coverFromRequest(r)
	if err != nil {
		s.renderBookError(w, r, existing, err.Error())
		return
	}
	if coverPath == "" {
		coverPath = existing.CoverPath
	}

	book := domain.Book{
		ID:        existing.ID,
		Title:     strings.TrimSpace(r.FormValue("title")),
		ISBN:      strings.TrimSpace(r.FormValue("isbn")),
		Year:      parseYear(r.FormValue("year")),
		Status:    domain.NormalizeStatus(r.FormValue("status")),
		CoverPath: coverPath,
		Notes:     strings.TrimSpace(r.FormValue("notes")),
		Authors:   authorsFromForm(r),
	}
	if _, err := s.store.UpdateBook(r.Context(), book); err != nil {
		s.renderBookError(w, r, book, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) personForm(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, http.StatusOK, "person_form", pageData{TitleKey: "title.add_person", CSRF: s.csrf(r)})
}

func (s *Server) peopleIndex(w http.ResponseWriter, r *http.Request) {
	people, err := s.store.ListPeople(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, http.StatusOK, "people_index", pageData{TitleKey: "title.people", CSRF: s.csrf(r), People: people})
}

func (s *Server) createPerson(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.render(w, r, http.StatusBadRequest, "person_form", pageData{TitleKey: "title.add_person", CSRF: s.csrf(r), Error: "error.invalid_form"})
		return
	}
	person := domain.Person{
		Name:     strings.TrimSpace(r.FormValue("name")),
		SortName: strings.TrimSpace(r.FormValue("sort_name")),
		Notes:    strings.TrimSpace(r.FormValue("notes")),
	}
	if _, err := s.store.CreatePerson(r.Context(), person); err != nil {
		s.render(w, r, http.StatusBadRequest, "person_form", pageData{TitleKey: "title.add_person", CSRF: s.csrf(r), Error: err.Error()})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) editPersonForm(w http.ResponseWriter, r *http.Request) {
	id := routeID(r, "id")
	person, err := s.store.GetPerson(r.Context(), id)
	if err != nil {
		s.renderNotFound(w, r, err)
		return
	}
	s.render(w, r, http.StatusOK, "person_form", pageData{
		TitleKey: "title.edit_person", CSRF: s.csrf(r), Person: person,
		Action: editPath("people", person.ID), Edit: true,
	})
}

func (s *Server) updatePerson(w http.ResponseWriter, r *http.Request) {
	id := routeID(r, "id")
	if err := r.ParseForm(); err != nil {
		s.render(w, r, http.StatusBadRequest, "person_form", pageData{TitleKey: "title.edit_person", CSRF: s.csrf(r), Error: "error.invalid_form"})
		return
	}
	person := domain.Person{
		ID:       id,
		Name:     strings.TrimSpace(r.FormValue("name")),
		SortName: strings.TrimSpace(r.FormValue("sort_name")),
		Notes:    strings.TrimSpace(r.FormValue("notes")),
	}
	if _, err := s.store.UpdatePerson(r.Context(), person); err != nil {
		s.render(w, r, http.StatusBadRequest, "person_form", pageData{
			TitleKey: "title.edit_person", CSRF: s.csrf(r), Error: err.Error(), Person: person,
			Action: editPath("people", id), Edit: true,
		})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) groupForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.formData(r, "title.add_group", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, http.StatusOK, "group_form", data)
}

func (s *Server) groupsIndex(w http.ResponseWriter, r *http.Request) {
	groups, err := s.store.ListGroups(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, http.StatusOK, "groups_index", pageData{TitleKey: "title.groups", CSRF: s.csrf(r), Groups: groups})
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderFormError(w, r, "group_form", "error.invalid_form")
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

func (s *Server) editGroupForm(w http.ResponseWriter, r *http.Request) {
	id := routeID(r, "id")
	group, err := s.store.GetGroup(r.Context(), id)
	if err != nil {
		s.renderNotFound(w, r, err)
		return
	}
	data, err := s.formData(r, "title.edit_group", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Group = group
	data.Action = editPath("groups", group.ID)
	data.Edit = true
	s.render(w, r, http.StatusOK, "group_form", data)
}

func (s *Server) updateGroup(w http.ResponseWriter, r *http.Request) {
	id := routeID(r, "id")
	if err := r.ParseForm(); err != nil {
		s.renderFormError(w, r, "group_form", "error.invalid_form")
		return
	}
	group := domain.Group{
		ID:       id,
		Name:     strings.TrimSpace(r.FormValue("name")),
		SortName: strings.TrimSpace(r.FormValue("sort_name")),
		Notes:    strings.TrimSpace(r.FormValue("notes")),
	}
	if _, err := s.store.UpdateGroup(r.Context(), group, r.Form["member_ids"]); err != nil {
		data, dataErr := s.formData(r, "title.edit_group", "")
		if dataErr != nil {
			http.Error(w, dataErr.Error(), http.StatusInternalServerError)
			return
		}
		group.Members = peopleFromIDs(r.Form["member_ids"])
		data.Group = group
		data.Action = editPath("groups", id)
		data.Edit = true
		data.Error = err.Error()
		s.render(w, r, http.StatusBadRequest, "group_form", data)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) formData(r *http.Request, titleKey, excludeMusicID string) (pageData, error) {
	people, err := s.store.ListPeople(r.Context())
	if err != nil {
		return pageData{}, err
	}
	groups, err := s.store.ListGroups(r.Context())
	if err != nil {
		return pageData{}, err
	}
	musicReferences, err := s.store.ListMusicReferences(r.Context(), "", excludeMusicID)
	if err != nil {
		return pageData{}, err
	}
	return pageData{TitleKey: titleKey, CSRF: s.csrf(r), People: people, Groups: groups, MusicReferences: musicReferences}, nil
}

func (s *Server) renderFormError(w http.ResponseWriter, r *http.Request, templateName, message string) {
	titleKey := map[string]string{
		"music_form": "title.add_music",
		"book_form":  "title.add_book",
		"group_form": "title.add_group",
	}[templateName]
	data, err := s.formData(r, titleKey, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Error = message
	s.render(w, r, http.StatusBadRequest, templateName, data)
}

func (s *Server) renderMusicError(w http.ResponseWriter, r *http.Request, item domain.MusicItem, message string) {
	data, err := s.formData(r, "title.edit_music", item.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Music = item
	data.Action = editPath("music", item.ID)
	data.Edit = item.ID != ""
	data.Error = message
	s.render(w, r, http.StatusBadRequest, "music_form", data)
}

func (s *Server) renderBookError(w http.ResponseWriter, r *http.Request, book domain.Book, message string) {
	data, err := s.formData(r, "title.edit_book", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Book = book
	data.Action = editPath("books", book.ID)
	data.Edit = book.ID != ""
	data.Error = message
	s.render(w, r, http.StatusBadRequest, "book_form", data)
}

func (s *Server) renderNotFound(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, translate(languageFromRequest(r), "error.not_found"), http.StatusNotFound)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
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

func authorsFromForm(r *http.Request) []domain.Person {
	authors := make([]domain.Person, 0)
	for _, id := range r.Form["person_ids"] {
		authors = append(authors, domain.Person{ID: id})
	}
	return authors
}

func peopleFromIDs(ids []string) []domain.Person {
	people := make([]domain.Person, 0, len(ids))
	for _, id := range ids {
		people = append(people, domain.Person{ID: id})
	}
	return people
}

func relatedFromForm(r *http.Request) []domain.MusicReference {
	related := make([]domain.MusicReference, 0, len(r.Form["related_ids"]))
	for _, id := range r.Form["related_ids"] {
		related = append(related, domain.MusicReference{ID: id})
	}
	return related
}

func (s *Server) musicSuggestions(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		s.render(w, r, http.StatusOK, "music_suggestions", pageData{})
		return
	}
	items, err := s.store.ListMusicReferences(r.Context(), query, r.URL.Query().Get("exclude"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, http.StatusOK, "music_suggestions", pageData{MusicReferences: items})
}

func routeID(r *http.Request, key string) string {
	value := chi.URLParam(r, key)
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

func editPath(kind, id string) string {
	return "/" + kind + "/" + url.PathEscape(id) + "/edit"
}
