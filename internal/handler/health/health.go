package health

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func New(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	// response data
	response := ResponseData{http.StatusOK, "health check ok"}
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}
