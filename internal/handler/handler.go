package handler

import (
	"database/sql"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/handler/health"
	"github.com/suburi-dev/gowiki/internal/session"
	"github.com/suburi-dev/gowiki/internal/handler/user"
	"net/http"
)

// wrapper 本体は /handler/health
func HealthHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := health.New(db, w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("HealthHandler failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// wrapper 本体は /handler/user
func UserHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := user.New(db, gs, w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("UserHandler failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
