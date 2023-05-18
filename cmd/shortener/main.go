package main

import (
	"github.com/SafronovRaff/URL-shortener/internal/controllers"
	"net/http"
)

//TODO: закоммить изменения и сделать мерж реквест(merge request), чтобы посмотреть работают ли автотесты
// TODO: вынести хендлеры отдельно в internal/controllers

func main() {
	http.HandleFunc("/", controllers.HandlerGet)
	http.HandleFunc("/{id}", controllers.HandlerPost)
	server := http.Server{
		Addr: "localhost:8080",
	}

	server.ListenAndServe()
}
