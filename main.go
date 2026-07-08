package main

import (
	"log"
	"net/http"

	"github.com/ValeraSun/PathEater/internal/network"
)

func main() {
	hub := network.NewHub()

	// websocket / сетевые хендлеры
	network.RegisterHandlers(hub)

	// раздача всего фронтенда из папки web
	http.Handle("/assets/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	http.Handle("/models/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	http.Handle("/textures/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	// корень сайта -> index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// если это запрос к статике, пробуем отдать файл
		if r.URL.Path != "/" {
			http.FileServer(http.Dir("./web")).ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, "./web/index.html")
	})

	log.Println("Server started on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}