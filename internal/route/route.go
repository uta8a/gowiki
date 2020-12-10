package route

import (
	"database/sql"
	"github.com/suburi-dev/gowiki/internal/handler"
	"net/http"
)

// Application State
type State struct {
	DB *sql.DB
}

func NewState(db *sql.DB) *State {
	return &State{DB: db}
}

func RegisterRoutes(db *sql.DB) {
	s := NewState(db)
	http.HandleFunc("/healthcheck", s.gen(handler.HealthHandler))
	http.HandleFunc("/users", s.gen(handler.UserHandler))
	// http.HandleFunc("/login", LoginHandler)
	// http.HandleFunc("/privatecheck", PrivateHandler)
}

// helper
type HandlerStateFunc func(db *sql.DB, w http.ResponseWriter, r *http.Request)

// convert f(db, w, r) -> db.f(w, r)
// non-local method is not allowed, so use f(db, w, r)
func (state *State) gen(handler HandlerStateFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(state.DB, w, r)
	}
}
