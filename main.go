package main

import (
	"log"

	"net/http"
	"github.com/ValeraSun/PathEater/internal/network"
)

func main() {

	hub := network.NewHub()

	network.RegisterHandlers(hub)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}