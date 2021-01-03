package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/suburi-dev/gowiki/internal/handler/article"
	"github.com/suburi-dev/gowiki/internal/handler/articleid"
	"github.com/suburi-dev/gowiki/internal/handler/group"
	"github.com/suburi-dev/gowiki/internal/handler/groupname"
	"github.com/suburi-dev/gowiki/internal/handler/health"
	"github.com/suburi-dev/gowiki/internal/handler/login"
	"github.com/suburi-dev/gowiki/internal/handler/private"
	"github.com/suburi-dev/gowiki/internal/handler/user"
	"github.com/suburi-dev/gowiki/internal/session"
)

// wrapper 本体は /handler/health
func HealthHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := health.New(db, w, r)
	if err != nil {
		log.Println("HealthHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// wrapper 本体は /handler/private
func PrivateHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := private.New(db, gs, w, r)
	if err != nil {
		log.Println("PrivateHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// wrapper 本体は /handler/user
func UserHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := user.New(db, gs, w, r)
	if err != nil {
		log.Println("UserHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// wrapper 本体は /handler/group
func GroupHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := group.New(db, gs, w, r)
	if err != nil {
		log.Println("GroupHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GroupNameHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := groupname.New(db, gs, w, r)
	if err != nil {
		log.Println("GroupNameHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// wrapper 本体は /handler/login
func LoginHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := login.New(db, gs, w, r)
	if err != nil {
		log.Println("LoginHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ArticleHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := article.New(db, gs, w, r)
	if err != nil {
		log.Println("ArticleHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ArticleIdHandler(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) {
	err := articleid.New(db, gs, w, r)
	if err != nil {
		log.Println("ArticleIdHandler failed: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
