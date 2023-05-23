package controllers

import (
	"fmt"
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

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    //  количество бит, которые будут использоваться для представления индекса символа
	letterIdxMask = 1<<letterIdxBits - 1 //  маска, используемая для извлечения битов индекса символа
	letterIdxMax  = 63 / letterIdxBits   // определяет максимальное количество индексов символов, которое помещается в 63 бита
)

var src = rand.NewSource(time.Now().UnixNano()) //используется для создания генератора случайных чисел на основе текущего времени
var urlMap = make(map[string]string)            //используется для хранения сокращенных URL на исходных URL
var mu sync.Mutex

func shorten(w http.ResponseWriter, r *http.Request) {
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
	keyURL := GenerateRandomString(10)
	mu.Lock()
	// Добавление значения URL в urlMap
	urlMap[keyURL] = urlString
	log.Printf("Добавлен URL в urlMap. Ключ: %s, Значение: %s", keyURL, urlString)
	mu.Unlock()
	// возвращаем сокращенный URL
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(keyURL))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid request method")
		return
	}
	key, ok := r.URL.Query()["link"]
	if ok {
		if _, ok := urlMap[key[0]]; !ok {
			genString := fmt.Sprint(rand.Int63n(1000))
			urlMap[genString] = key[0]
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusAccepted)
			linkString := fmt.Sprintf("<a href=\"http://localhost:8080/short/%s\">http://localhost:8080/short/%s</a>", genString, genString)
			fmt.Fprintf(w, "Added shortlink\n")
			fmt.Fprintf(w, linkString)
			return
		}
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Already have this link")
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Failed to add link")
	return
}
func Increase(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathArgs := strings.Split(path, "/")
	log.Printf("Redirected to: %s", urlMap[pathArgs[2]])

	http.Redirect(w, r, urlMap[pathArgs[2]], http.StatusPermanentRedirect)

	return
}

func increase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid request method")
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
	mu.Lock()
	originalURL, ok := urlMap[decodedURL]
	log.Printf("Извлечен URL из urlMap. Ключ: %s, Значение: %s, Найден: %v", decodedURL, originalURL, ok)
	mu.Unlock()

	if !ok {
		http.Error(w, "недопустимый идентификатор URL-адреса", http.StatusNotFound)
		return
	}

	// Возвращаем оригинальный URL
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Функция генерирует случайную строку длиной "n" из  байтового слайса
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
