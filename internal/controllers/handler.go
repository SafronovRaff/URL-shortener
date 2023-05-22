package controllers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())
var urlMap = make(map[string]string)

func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	resp, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	short := string(resp)
	shortURL := shortenURL(string(resp))
	urlMap[short] = short

	w.Header().Set("content-type", "text/plain;charset=utf-8 ")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func Increase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	//id := r.URL.Path[len("/"):]
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}

	url, ok := urlMap[id]
	if !ok {
		http.Error(w, "invalid URL ID", http.StatusBadRequest)
		return
	}

	//http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Функция для генерации сокращенного URL
func shortenURL(url string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(15) - 5
	log.Printf("рандом длинны строки %d", n)
	res := randStringBytes(n)
	log.Printf("рандом строка %s", res)
	url = res
	return url
}

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
