package controllers

import (
	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var urlmap = make(map[string]string)

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
	defer r.Body.Close()
	// Преобразуем данные в строку URL
	/*urlString, err := url.PathUnescape(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/
	urlString := string(b)
	// Вывод значения URL в лог
	log.Printf("Извлеченное значение URL: %s", urlString)

	// Генерация случайной строки в качестве ключа
	keyURL := maintenance.GenerateRandomString(10)

	// Добавление значения URL в urlMap
	//savedURL := maintenance.NewMap().Add(keyURL, urlString)

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

	vars := mux.Vars(r)

	id := vars["id"]

	if r.Method != http.MethodGet {
		http.Error(w, "Недопустимый метод запроса", http.StatusBadRequest)
		return
	}

	log.Printf("URL.Path: %s", id)

	originalURL, ok := urlmap[id]
	if !ok {
		http.Error(w, "originalURL не найден", http.StatusBadRequest)
	}
	log.Printf("найден originalURL: %s", originalURL)

	// Возвращаем оригинальный URL
	w.Header().Set("Location", originalURL)
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)

}
