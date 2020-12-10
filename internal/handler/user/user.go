package user

import (
	// "encoding/json"
	// "fmt"
	"database/sql"
	"github.com/suburi-dev/gowiki/internal/auth"
  "net/http"
  "golang.org/x/crypto/bcrypt"
)

func New(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	// POST
	if r.Method == http.MethodPost {
		// signup
		err := signup(db, w, r)
		if err != nil {
			return err
		}
		return nil
	}
	http.NotFound(w, r)
	return nil
}

// signup
func signup(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
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
  // username identity check
  var exists bool
  query := fmt.Sprintf("SELECT EXISTS (SELECT username FROM users WHERE username = %s)", username)
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
  query = fmt.Sprintf("INSERT INTO users(username, password_hash) VALUES($1,$2)", username, hashed)
  _, err = db.Exec(query)
  if err != nil {
    return err
  }
  // Session

  return nil
}
