package main

import (
	"github.com/SafronovRaff/URL-shortener/internal/controllers"
	"github.com/SafronovRaff/URL-shortener/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	storage.NewMap() //создаём мапу

	router := mux.NewRouter()
	router.HandleFunc("/", controllers.Shorten).Methods(http.MethodPost)
	router.HandleFunc("/{id}", controllers.Increase).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))

}
