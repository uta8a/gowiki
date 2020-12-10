package route

import (
	"database/sql"
	"github.com/suburi-dev/gowiki/internal/session"
	"github.com/suburi-dev/gowiki/internal/handler"
	"net/http"
)

// Application State
type State struct {
  DB *sql.DB
  Session *session.Manager
}

func NewState(db *sql.DB, gs *session.Manager) *State {
	return &State{DB: db, Session: gs}
}

func RegisterRoutes(db *sql.DB, gs *session.Manager) {
	s := NewState(db, gs)
	http.HandleFunc("/healthcheck", s.gen(handler.HealthHandler))
	http.HandleFunc("/users", s.gen(handler.UserHandler))
	// http.HandleFunc("/login", LoginHandler)
	// http.HandleFunc("/privatecheck", PrivateHandler)
}

// helper
type HandlerStateFunc func(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request)

// convert f(db, w, r) -> db.f(w, r)
// non-local method is not allowed, so use f(db, w, r)
func (state *State) gen(handler HandlerStateFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(state.DB, state.Session, w, r)
	}
}
