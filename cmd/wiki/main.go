package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type ResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Application State
type State struct {
	db *sql.DB
}

func InitEnv(k string) string {
	// k: key v: value
	v, ok := os.LookupEnv(k)
	// if unset env key:value, Fatal
	if !ok {
		log.Fatal("Environment value is not set, key: ", k)
	}
	return v
}
func NewState(db *sql.DB) *State {
	return &State{db: db}
}
func Health() {}

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

func RegisterRoutes(db *sql.DB) {
	_ = NewState(db)
	http.HandleFunc("/healthcheck", HealthHandler)
  // http.HandleFunc("/privatecheck", privateHandler)
	// http.HandleFunc("/users", userHandler)
	// http.HandleFunc("/login", loginHandler)
}

func main() {
	// db init
	dbhost := InitEnv("POSTGRES_HOST")
	dbport := InitEnv("POSTGRES_PORT")
	dbuser := InitEnv("POSTGRES_USER")
	dbpass := InitEnv("POSTGRES_PASSWORD")
	dbname := InitEnv("POSTGRES_DB")
	dburl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)
	log.Println("DB URL: ", dburl)
	db, err := sql.Open("postgres", dburl) // return *sql.DB, error
	if err != nil {
		log.Fatal("DB connection failed: ", err)
  }
  
  // Routes
	RegisterRoutes(db)

  // Server start
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
