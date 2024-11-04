package main

import (
	"caibai/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
