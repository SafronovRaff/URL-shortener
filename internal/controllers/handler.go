package controllers

import (
	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
	"io"
	"log"
	"net/http"
	"net/url"
)

var urlmap map[string]string

func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Недопустимый метод запроса", http.StatusBadRequest)
		return
	}
	// Считываем данные из тела запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразуем данные в строку URL
	urlString, err := url.PathUnescape(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вывод значения URL в лог
	log.Printf("Извлеченное значение URL: %s", urlString)

	// Генерация случайной строки в качестве ключа
	keyURL := maintenance.GenerateRandomString(10)

	// Добавление значения URL в urlMap
	//savedURL := maintenance.NewMap().Add(keyURL, urlString)
	urlmap = make(map[string]string)

	urlmap[keyURL] = urlString
	log.Printf("Добавлен URL в urlMap. Ключ: %s, Значение: %s", keyURL, urlString)

	// возвращаем сокращенный URL
	//w.Header().Set("content-type", "text/plain; charset=utf-8")
	//w.WriteHeader(http.StatusCreated)
	//
	//_, err = w.Write([]byte(keyURL))
	//if err != nil {
	//	http.Error(w, err.Error(), 400)
	//	return
	//}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + keyURL))

}

func Increase(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Недопустимый метод запроса", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[1:]
	log.Printf("URL.Path: %s", id)

	originalURL, ok := urlmap[id]
	if !ok {
		http.Error(w, "originalURL не найден", http.StatusBadRequest)
	}
	log.Printf("найден originalURL: %s", originalURL)

	// Возвращаем оригинальный URL
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
