package route

import (
	"database/sql"
	"github.com/suburi-dev/gowiki/internal/handler"
	"github.com/suburi-dev/gowiki/internal/session"
	"net/http"
)

// Application State
type State struct {
	DB      *sql.DB
	Session *session.Manager
}

func NewState(db *sql.DB, gs *session.Manager) *State {
	return &State{DB: db, Session: gs}
}

func RegisterRoutes(db *sql.DB, gs *session.Manager) {
	s := NewState(db, gs)
	http.HandleFunc("/healthcheck", s.gen(handler.HealthHandler))
	http.HandleFunc("/privatecheck", s.gen(handler.PrivateHandler))
	http.HandleFunc("/users", s.gen(handler.UserHandler))
	http.HandleFunc("/groups", s.gen(handler.GroupHandler))
	http.HandleFunc("/groups/", s.gen(handler.GroupNameHandler))
	http.HandleFunc("/login", s.gen(handler.LoginHandler))
	http.HandleFunc("/articles", s.gen(handler.ArticleHandler))
	http.HandleFunc("/articles/", s.gen(handler.ArticleIdHandler))
}

// helper
type HandlerStateFunc func(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request)

// convert f(db, gs, w, r) -> state.f(w, r)
// non-local method is not allowed, so use f(db, gs, w, r)
func (state *State) gen(handler HandlerStateFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(state.DB, state.Session, w, r)
	}
}
