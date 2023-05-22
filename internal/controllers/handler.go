package controllers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
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

	/*	//проверяем, что метод запроса является POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
			return
		}*/
	// считываем данные из тела запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//генерируем сокращенный URL с помощью функции shortenURL и закидываем в мапу где сокрURl ключ, а URL значение.
	short := string(b)
	log.Printf("URL значение - %s", short)
	shortURL := shortenURL()
	mu.Lock()
	urlMap[shortURL] = short
	mu.Unlock()
	// возвращаем сокращенный URL
	w.Header().Set("content-type", "text/plain;charset=utf-8 ")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func Increase(w http.ResponseWriter, r *http.Request) {
	/*//проверяем, что метод запроса является Get
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}*/
	//считываем id

	id := chi.URLParam(r, "id")
	log.Printf("id- %s", id)
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}
	parsedURL, err := url.Parse(id)
	if err != nil || parsedURL.Scheme == "" {
		http.Error(w, "invalid URL format", http.StatusBadRequest)
		return
	}

	//ищем в мапе оригинальный URL
	mu.Lock()
	url, ok := urlMap[id]
	log.Printf("url %s", url)
	mu.Unlock()
	if !ok {
		http.Error(w, "invalid URL ID", http.StatusBadRequest)
		return
	}

	log.Printf("оригинальный URL %s", url)
	//возвращаем оригинальный URL
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Функция для генерации сокращенного URL
func shortenURL() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(15) + 5
	res := randStringBytes(n)
	log.Printf("рандом строка %s", res)
	return res
}

// Функция генерирует случайную строку длиной "n" из  байтового слайса
func randStringBytes(n int) string {
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
