package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/suburi-dev/gowiki/internal/env"
	"github.com/suburi-dev/gowiki/internal/route"
	"log"
	"net/http"
)

func main() {
	// db init
	dbhost := env.Init("POSTGRES_HOST")
	dbport := env.Init("POSTGRES_PORT")
	dbuser := env.Init("POSTGRES_USER")
	dbpass := env.Init("POSTGRES_PASSWORD")
	dbname := env.Init("POSTGRES_DB")
	dburl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)
	log.Println("DB URL: ", dburl)
	db, err := sql.Open("postgres", dburl) // return *sql.DB, error
	if err != nil {
		log.Fatal("DB Open failed: ", err)
  }
  err = db.Ping()
  if err != nil {
    log.Fatal("DB Ping failed: ", err)
  }

	// Routes
	route.RegisterRoutes(db)

	// Server start
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
