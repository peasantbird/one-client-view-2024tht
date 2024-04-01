package main

import (
	"golang-api/pkg/api"
	"golang-api/pkg/config"
	"golang-api/pkg/db"
	"log"
	"net/http"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		panic(err)
	}

	router := api.Router(&api.Handler{DB: db})

	// Use the router as the default HTTP handler
	http.Handle("/", router)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":"+config.GetPort(), router))
}
