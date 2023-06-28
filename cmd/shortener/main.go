package main

import (
	"github.com/SafronovRaff/URL-shortener/internal/controllers"
	"github.com/SafronovRaff/URL-shortener/internal/generate"
	"github.com/SafronovRaff/URL-shortener/internal/service"
	"github.com/SafronovRaff/URL-shortener/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// newGenerate -генератор рандомной строки
	newGenerate := generate.NewGenerate()
	//newStorage - БД/Мапа
	newStorage := storage.NewStorage()
	newService := service.NewService(newGenerate, newStorage)
	handlers := controllers.NewHandlers(newService)

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Shortened).Methods(http.MethodPost)
	router.HandleFunc("/{id}", handlers.Increase).Methods(http.MethodGet)

	router.HandleFunc("/api/shorten", handlers.ShortenHandler).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))

}
