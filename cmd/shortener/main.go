package main

import (
	"github.com/SafronovRaff/URL-shortener/internal/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	//maintenance.NewMap() //создаём мапу
	router := mux.NewRouter()

	router.HandleFunc("/", controllers.Shorten).Methods(http.MethodPost)
	router.HandleFunc("/{id}", controllers.Increase).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)

}

//import (
//	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
//	"github.com/gorilla/mux"
//	"io"
//	"log"
//	"net/http"
//)
//
//var urls = make(map[string]string)
//
//func ShortURL(w http.ResponseWriter, r *http.Request) {
//	b, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//	}
//	defer r.Body.Close()
//
//	url := string(b)
//	//key := genShortenURL(url)
//	key := maintenance.GenerateRandomString(10)
//
//	urls[key] = url
//
//	log.Println(urls)
//	w.Header().Set("content-type", "text/plain; charset=utf-8")
//	w.WriteHeader(http.StatusCreated)
//	w.Write([]byte("http://localhost:8080/" + key))
//}
//
//func GetID(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id := vars["id"]
//	log.Println(id)
//	if v, ok := urls[id]; ok {
//		w.Header().Set("Location", v)
//		http.Redirect(w, r, v, http.StatusTemporaryRedirect)
//		return
//	} else {
//		w.WriteHeader(http.StatusBadRequest)
//	}
//}
//
//func main() {
//	r := mux.NewRouter()
//	r.HandleFunc("/", ShortURL).Methods(http.MethodPost)
//	r.HandleFunc("/{id}", GetID).Methods(http.MethodGet)
//	log.Fatal(http.ListenAndServe(":8080", r))
//}
