package web

import (
	"crypto/subtle"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const sessionCookie = "albumer_session"

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, http.StatusOK, "login", pageData{TitleKey: "title.login"})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.render(w, r, http.StatusBadRequest, "login", pageData{TitleKey: "title.login", Error: "error.invalid_form"})
		return
	}
	if !s.passwordMatches(r.FormValue("password")) {
		s.render(w, r, http.StatusUnauthorized, "login", pageData{TitleKey: "title.login", Error: "error.invalid_password"})
		return
	}

	id := newToken()
	s.mu.Lock()
	s.sessions[id] = session{CSRF: newToken()}
	s.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookie)
	if err == nil {
		s.mu.Lock()
		delete(s.sessions, cookie.Value)
		s.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := s.currentSession(r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := s.currentSession(r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (s *Server) csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		sess, ok := s.currentSession(r)
		if !ok || subtle.ConstantTimeCompare([]byte(sess.CSRF), []byte(r.FormValue("csrf"))) != 1 {
			http.Error(w, translate(languageFromRequest(r), "error.invalid_csrf"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) currentSession(r *http.Request) (session, bool) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return session{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[cookie.Value]
	return sess, ok
}

func (s *Server) passwordMatches(password string) bool {
	if s.cfg.AdminPasswordHash == "" {
		return subtle.ConstantTimeCompare([]byte(password), []byte("admin")) == 1
	}
	if err := bcrypt.CompareHashAndPassword([]byte(s.cfg.AdminPasswordHash), []byte(password)); err == nil {
		return true
	}
	return subtle.ConstantTimeCompare([]byte(password), []byte(s.cfg.AdminPasswordHash)) == 1
}
