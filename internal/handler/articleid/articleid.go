package articleid

import (
  "database/sql"
	"github.com/suburi-dev/gowiki/internal/session"
  "net/http"
)
func New(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
  return nil
}
