package main

import (
	"log"
	"net/http"

	"github.com/ValeraSun/PathEater/internal/network"
)

func main() {
	hub := network.NewHub()

	network.RegisterHandlers(hub)

	port := ":8080"
	log.Printf("Server starting on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}