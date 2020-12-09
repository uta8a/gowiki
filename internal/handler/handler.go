package handler

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type ResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// response data
	response := ResponseData{http.StatusOK, "health check ok"}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
}
