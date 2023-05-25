package controllers

import (
	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
	"io"
	"log"
	"net/http"
)

var urlmap map[string]string

func Shorten(w http.ResponseWriter, r *http.Request) {

	// считываем данные из тела запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	urlString := string(b)
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
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(keyURL))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func Increase(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path
	log.Printf("URL.Path: %s", id)
	id = id[1:]

	originalURL, ok := urlmap[id]
	log.Printf("найден originalURL: %s", originalURL)
	if !ok {
		http.Error(w, "originalURL не найден", http.StatusBadRequest)
	}
	// Возвращаем оригинальный URL
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)

}
