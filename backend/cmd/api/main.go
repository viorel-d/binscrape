package main

import (
	"log"
	"net/http"

	"github.com/viorel-d/binscrape/backend/config"
	"github.com/viorel-d/binscrape/backend/handlers"
)

func main() {
	http.HandleFunc("/list", handlers.GetAllPasteBinItemsHandler)
	http.HandleFunc("/details/{id}", handlers.GetPasteBinItemDetailsHandler)
	log.Printf("Starting server at address: %v\n", config.ServerAddress)
	http.ListenAndServe(config.ServerAddress, nil)
}
