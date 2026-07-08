package main

import (
	"log"
	"net/http"
	"server/internal/infrastructure/network"
)

func main() {
	network.RegisterHandlers()
	
	http.Handle("/", http.FileServer(http.Dir("./static")))
	
	port := ":8080"
	log.Println("Сервер запущен на http://localhost", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic(err)
	}
}