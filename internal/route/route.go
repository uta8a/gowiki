package route

import (
	"database/sql"
	"github.com/suburi-dev/gowiki/internal/handler"
	"net/http"
)

// Application State
type State struct {
	db *sql.DB
}

func NewState(db *sql.DB) *State {
	return &State{db: db}
}

func RegisterRoutes(db *sql.DB) {
	_ = NewState(db)
	http.HandleFunc("/healthcheck", handler.HealthHandler)
	// http.HandleFunc("/privatecheck", privateHandler)
	// http.HandleFunc("/users", userHandler)
	// http.HandleFunc("/login", loginHandler)
}
