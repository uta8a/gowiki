package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/auth"
	"github.com/suburi-dev/gowiki/internal/session"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"time"
)

type ResponseData struct {
	Status   int    `json:"status"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func New(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	// POST
	if r.Method == http.MethodPost {
		// signup
		err := signup(db, gs, w, r)
		if err != nil {
			return err
		}
		return nil
	}
	http.NotFound(w, r)
	return nil
}

// signup
func signup(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// validate
	username := r.FormValue("username")
	password := r.FormValue("password")
	err = auth.Validate(username, password)
	if err != nil {
		return err
	}
	log.Printf("validate ok: %s", username)
	// username identity check
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT username FROM users WHERE username = '%s')", username)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("username %s already exists", username)
	}
	// hash
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	hashed := string(hashedBytes)
	// db insert
	query = fmt.Sprintf("INSERT INTO users(username, password_hash) VALUES('%s','%s')", username, hashed)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	// Session
	sess := gs.SessionStart(w, r)
	// cookie
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0) // TODO fix
	cookie := http.Cookie{Name: "SESSIONID", Value: sess.SessionID(), Expires: expiration}
	http.SetCookie(w, &cookie)
	// response data
	response := ResponseData{http.StatusOK, username, "signup ok"} // TODO fix status code to 201
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}
