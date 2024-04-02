package main

import (
	"golang-api/internal/api"
	"golang-api/internal/config"
	"golang-api/internal/db"
	"log"
	"net/http"
)

func main() {
	config := config.NewConfig()
	config.Load()

	db, err := db.Connect(config)
	if err != nil {
		panic(err)
	}

	router := api.Router(&api.Handler{DB: db})

	// Use the router as the default HTTP handler
	http.Handle("/", router)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
