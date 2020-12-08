package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
)
type ResponseData struct {
  Status int `json:"status"`
  Message string `json:"message"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
  http.HandleFunc("/healthcheck", healthHandler)
  // http.HandleFunc("/privatecheck", privateHandler)
  // http.HandleFunc("/users", userHandler)
  // http.HandleFunc("/login", loginHandler)

  err := http.ListenAndServe(":9000", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
