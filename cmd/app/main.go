package main

import (
	"log"
	"net/http"

	"github.com/ValeraSun/PathEater/internal/network"
)

func main() {
	network.RegisterHandlers()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Println("Сервер слушает :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}