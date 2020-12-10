package private

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/session"
	"net/http"
)

type ResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func New(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	// middleware auth
	// sess, err := SessionCheck(gs)
	ok := gs.SessionCheck(w, r)
	if !ok {
		http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
	}
	sess := gs.SessionStart(w, r)
	username := sess.Get("username")
	// response data
	response := ResponseData{http.StatusOK, fmt.Sprintf("Welcome %s!\nThis is Private Page.", username)}
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}
