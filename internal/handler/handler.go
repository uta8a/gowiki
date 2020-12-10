package handler

import (
	"database/sql"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/handler/health"
	"github.com/suburi-dev/gowiki/internal/handler/private"
	"github.com/suburi-dev/gowiki/internal/handler/user"
	"github.com/suburi-dev/gowiki/internal/session"
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

// wrapper 本体は /handler/private
func PrivateHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := private.New(db, gs, w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("PrivateHandler failed: %s", err.Error()), http.StatusInternalServerError)
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
