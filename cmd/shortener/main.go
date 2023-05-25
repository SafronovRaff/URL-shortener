package main

import (
	"github.com/SafronovRaff/URL-shortener/internal/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

//TODO: закоммить изменения и сделать мерж реквест(merge request), чтобы посмотреть работают ли автотесты
// TODO: вынести хендлеры отдельно в internal/controllers

func main() {

	//maintenance.NewMap() //создаём мапу
	router := mux.NewRouter()

	router.HandleFunc("/", controllers.Shorten)
	router.HandleFunc("/{id}", controllers.Increase)

	http.ListenAndServe(":8080", router)

}
