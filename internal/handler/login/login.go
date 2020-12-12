package login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/auth"
	"github.com/suburi-dev/gowiki/internal/session"
  "golang.org/x/crypto/bcrypt"
  "log"
  "net/http"
  "io/ioutil"
)

type ResponseData struct {
	Status   int    `json:"status"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
type DelRes struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
}
type User struct {
  Username string `json:"username"`
  Password string `json:"password"`
}
func New(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	// POST
	if r.Method == http.MethodPost {
		err := loginCheck(db, gs, w, r)
		if err != nil {
			return err
		}
		return nil
  }
  // DELETE
	if r.Method == http.MethodDelete {
		err := logout(db, gs, w, r)
		if err != nil {
			return err
		}
		return nil
	}
	http.NotFound(w, r)
	return nil
}

// login
func loginCheck(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    return err
  }
  var u User
  if err := json.Unmarshal(body, &u); err != nil {
    return err
  }
	// validate
	username := u.Username
	password := u.Password
	err = auth.Validate(username, password)
	if err != nil {
		return err
	}
	log.Printf("validate ok: %s", username)
	// get password hash
	var passwordHash string
	query := fmt.Sprintf("SELECT password_hash FROM users WHERE username = '%s'", username)
	err = db.QueryRow(query).Scan(&passwordHash)
	if err != nil {
		return err
	}
	// compare hash & password
  err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return err
  }
	// SessionStart
	sess := gs.SessionStart(w, r)
	sess.Set("username", username)
	// response data
	response := ResponseData{http.StatusOK, username, "login ok"} // TODO fix status code to 201
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}

func logout(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
  ok := gs.SessionCheck(w, r)
  if !ok {
    http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
  }
  // get username from session then destroy session
  sess := gs.SessionStart(w, r)
  username := sess.Get("username")
  gs.SessionDestroy(w, r)
  // response data
  response := DelRes{http.StatusOK, fmt.Sprintf("%s Logout Success", username)}
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}
