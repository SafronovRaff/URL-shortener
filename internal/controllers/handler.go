package controllers

import (
	"io"
	"net/http"
	"strings"
)

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
	id := strings.TrimPrefix(r.URL.Path, "/")
	url, ok := urlMap[id]
	if !ok {
		http.Error(w, "invalid URL ID", http.StatusBadRequest)
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Функция для генерации сокращенного URL
func shortenURL(url string) string {
	//TODO тут  будет алгоритм генерации сокращенного URL
	return url
}
