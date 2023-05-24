package controllers

import (
	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano()) //используется для создания генератора случайных чисел на основе текущего времени
var urlMap = make(map[string]string)            //используется для хранения сокращенных URL на исходных URL
var mu sync.Mutex

func Shorten(w http.ResponseWriter, r *http.Request) {
	// считываем данные из тела запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	urlString, err := url.QueryUnescape(string(b))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// Вывод значения URL в лог
	log.Printf("Извлеченное значение URL: %s", urlString)
	// Генерация случайной строки в качестве ключа
	keyURL := maintenance.GenerateRandomString(10)

	// Добавление значения URL в urlMap
	maintenance.NewMap().Add(keyURL, urlString)

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
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	id := chi.URLParam(r, "id")
	log.Printf("id- %s", id)
	if id == "" {
		http.Error(w, "параметр id пуст", http.StatusBadRequest)
		return
	}

	// Добавляем схему протокола, если она отсутствует
	if !strings.HasPrefix(id, "http://") && !strings.HasPrefix(id, "https://") {
		id = "http://" + id
	}

	parsedURL, err := url.Parse(id)
	if err != nil {
		http.Error(w, "недопустимый формат URL-адреса", http.StatusBadRequest)
		return
	}

	// Декодируем URL
	decodedURL := parsedURL.String()
	decodedURL, err = url.PathUnescape(decodedURL)
	if err != nil {
		http.Error(w, "ошибка декодирования URL-адреса", http.StatusBadRequest)
		return
	}

	// Ищем оригинальный URL в urlMap

	originalURL, ok := maintenance.NewMap().Get(decodedURL)
	log.Printf("Извлечен URL из urlMap. Ключ: %s, Значение: %s, Найден: %v", decodedURL, originalURL, ok)
	if ok != nil {
		http.Error(w, "недопустимый идентификатор URL-адреса", http.StatusNotFound)
		return
	}

	// Возвращаем оригинальный URL
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
