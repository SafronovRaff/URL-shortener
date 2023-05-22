package controllers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
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

func Shorten(w http.ResponseWriter, r *http.Request) {

	//проверяем, что метод запроса является POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	// считываем данные из тела запроса
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//генерируем сокращенный URL с помощью функции shortenURL и закидываем в мапу где сокрURl ключ, а URL значение.
	short := string(b)
	shortURL := shortenURL(string(b))
	urlMap[shortURL] = short
	// возвращаем сокращенный URL
	w.Header().Set("content-type", "text/plain;charset=utf-8 ")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func Increase(w http.ResponseWriter, r *http.Request) {
	//проверяем, что метод запроса является Get
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	//считываем id
	//id := r.URL.Path[len("/"):]
	//id := chi.URLParam(r, "id")
	id := chi.URLParam(r, "/")
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}
	//ищем в мапе оригинальный URL
	url, ok := urlMap[id]
	if !ok {
		http.Error(w, "invalid URL ID", http.StatusBadRequest)
		return
	}
	//возвращаем оригинальный URL
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Функция для генерации сокращенного URL
func shortenURL(url string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(15) + 5
	//	log.Printf("рандом длинны строки %d", n)
	res := randStringBytes(n)
	//	log.Printf("рандом строка %s", res)
	url = res
	return url
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
