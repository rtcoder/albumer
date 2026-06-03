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
		protected.Get("/books/new", s.bookForm)
		protected.Post("/books/new", s.createBook)
		protected.Get("/people/new", s.personForm)
		protected.Post("/people/new", s.createPerson)
		protected.Get("/groups/new", s.groupForm)
		protected.Post("/groups/new", s.createGroup)
	})

	return r
}

func (s *Server) render(w http.ResponseWriter, status int, name string, data any) {
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
