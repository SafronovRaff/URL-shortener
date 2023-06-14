package controllers

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type services interface {
	ShortenedURL(b string) (string, error)
	IncreaseURL(id string) (string, error)
}

type Handlers struct {
	services services
}

func NewHandlers(services services) *Handlers {
	return &Handlers{services: services}
}

func (h *Handlers) Shortened(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "недопустимый метод запроса", http.StatusBadRequest)
		return
	}
	// Считываем данные из тела запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keyURL, err := h.services.ShortenedURL(string(b))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("http://localhost:8080/" + keyURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handlers) Increase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Недопустимый метод запроса", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]

	originalURL, err := h.services.IncreaseURL(id)
	if err == nil { // Возвращаем оригинальный URL
		log.Printf("найден originalURL: %s", originalURL)
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		http.Error(w, "originalURL не найден", http.StatusBadRequest)
	}

}
