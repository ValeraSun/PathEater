package main

import (
	"log"
	"net/http"

	"github.com/ValeraSun/PathEater/internal/network"
)

func main() {

	network.RegisterHandlers()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := ":8080"
	log.Println("Сервер запущен на http://localhost", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic(err)
	}
}
