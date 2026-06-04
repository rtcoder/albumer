package web

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"albumer/internal/config"
	"albumer/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg       config.Config
	store     storage.Store
	templates *template.Template
	sessions  map[string]session
	mu        sync.RWMutex
}

type session struct {
	CSRF string
}

func NewServer(cfg config.Config, store storage.Store) *Server {
	return &Server{
		cfg:       cfg,
		store:     store,
		templates: parseTemplates(),
		sessions:  map[string]session{},
	}
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/login", s.loginForm)
	r.Post("/login", s.login)
	r.Get("/language/{lang}", s.setLanguage)
	r.Get("/theme/{theme}", s.setTheme)
	r.Get("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.cfg.UploadDir))).ServeHTTP)
	r.Get("/favicon/*", http.StripPrefix("/favicon/", http.FileServer(http.Dir("favicon"))).ServeHTTP)
	r.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS()))).ServeHTTP)

	r.Group(func(protected chi.Router) {
		protected.Use(s.authMiddleware)
		protected.Use(s.csrfMiddleware)
		protected.Get("/", s.catalog)
		protected.Post("/logout", s.logout)
		protected.Get("/items", s.catalog)
		protected.Get("/music/new", s.musicForm)
		protected.Post("/music/new", s.createMusic)
		protected.Get("/music/suggestions", s.musicSuggestions)
		protected.Get("/music/{id}/edit", s.editMusicForm)
		protected.Post("/music/{id}/edit", s.updateMusic)
		protected.Get("/books/new", s.bookForm)
		protected.Post("/books/new", s.createBook)
		protected.Get("/books/{id}/edit", s.editBookForm)
		protected.Post("/books/{id}/edit", s.updateBook)
		protected.Get("/people", s.peopleIndex)
		protected.Get("/people/new", s.personForm)
		protected.Post("/people/new", s.createPerson)
		protected.Get("/people/{id}/edit", s.editPersonForm)
		protected.Post("/people/{id}/edit", s.updatePerson)
		protected.Get("/groups", s.groupsIndex)
		protected.Get("/groups/new", s.groupForm)
		protected.Post("/groups/new", s.createGroup)
		protected.Get("/groups/{id}/edit", s.editGroupForm)
		protected.Post("/groups/{id}/edit", s.updateGroup)
	})

	return r
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, status int, name string, data pageData) {
	data.Lang = languageFromRequest(r)
	data.Theme = themeFromRequest(r)
	data.NextURL = r.URL.RequestURI()
	if data.TitleKey != "" {
		data.Title = data.Tr(data.TitleKey)
	}
	if data.Error != "" {
		data.Error = translateError(data.Lang, data.Error)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := s.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newToken() string {
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes[:])
}

func uploadPath(uploadDir, name string) string {
	return filepath.Join(uploadDir, name)
}
