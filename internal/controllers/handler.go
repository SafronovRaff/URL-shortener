package controllers

import (
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
	keyURL := GenerateRandomString(10)

	// Добавление значения URL в urlMap
	mu.Lock()
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

func Increase(w http.ResponseWriter, r *http.Request) {
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
	decodedURL, err := url.PathUnescape(parsedURL.String())
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
