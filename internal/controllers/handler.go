package controllers

import (
	"encoding/json"
	"net/http"
)

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	q := r.URL.Query().Get("query")
	if q == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func HandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[1:]
	if id == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}

	originalURL := "http://example.com/original-url"

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
