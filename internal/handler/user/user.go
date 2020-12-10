package user

import (
	// "encoding/json"
	// "fmt"
	"database/sql"
	"github.com/suburi-dev/gowiki/internal/state"
	"net/http"
)

func New(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
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

}
