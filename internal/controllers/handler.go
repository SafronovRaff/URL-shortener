package controllers

import (
	"encoding/json"
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

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
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

func (h *Handlers) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "недопустимый метод запроса", http.StatusBadRequest)
		return
	}
	// Проверяем Content-Type заголовок на правильное значение
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Недопустимый Content-Type", http.StatusBadRequest)
		return
	}
	// Читаем тело запроса
	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Не верное тело запроса", http.StatusBadRequest)
		return
	}
	shortURL, err := h.services.ShortenedURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Формируем ответ в виде JSON
	resp := ShortenResponse{Result: shortURL}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Не удалось организовать ответ", http.StatusInternalServerError)
		return
	}
	// Устанавливаем правильный Content-Type заголовок
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
