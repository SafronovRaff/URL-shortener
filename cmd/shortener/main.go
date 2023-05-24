package main

import (
	"github.com/SafronovRaff/URL-shortener/internal/controllers"
	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
	"net/http"
)

//TODO: закоммить изменения и сделать мерж реквест(merge request), чтобы посмотреть работают ли автотесты
// TODO: вынести хендлеры отдельно в internal/controllers

func main() {

	maintenance.NewMap() //создаём мапу

	http.HandleFunc("/", controllers.Shorten)
	http.HandleFunc("/{id}", controllers.Increase)
	server := http.Server{
		Addr: "localhost:8080",
	}

	server.ListenAndServe()
}
