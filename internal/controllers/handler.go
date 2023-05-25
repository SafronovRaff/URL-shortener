package controllers

import (
	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Result struct {
	Link   string
	Code   string
	Status int
}

var urlmap map[string]string

func isValidUrl(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}
func Shorten(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "<h1>hellow</h1>")
	//templ, _ := template.ParseFiles("templates/index.html")
	//result := Result{}
	//if r.Method == "POST" {
	//	if !isValidUrl(r.FormValue("s")) {
	//		log.Printf("Что-то не так")
	//		result.Status = 400
	//		result.Link = ""
	//	} else {
	//		result.Link = r.FormValue("s")
	//		result.Code = maintenance.GenerateRandomString(10)
	//		maintenance.NewMap().Add(result.Link, result.Code)
	//		log.Printf("Добавлен URL в urlMap. Ключ: %s, Значение: %s", result.Link, result.Code)
	//		result.Status = http.StatusCreated
	//	}
	//}
	//templ.Execute(w, result)
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
	//savedURL := maintenance.NewMap().Add(keyURL, urlString)
	urlmap = make(map[string]string)
	urlmap[keyURL] = urlString
	log.Printf("Добавлен URL в urlMap. Ключ: %s, Значение: %s", keyURL, urlString)

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

	//	id := mux.Vars(r)
	//
	//originalURL, err := maintenance.NewMap().Get(id["key"])
	////log.Printf("Извлечен URL из urlMap. Ключ: %s, Значение: %s, Найден: %v", decodedURL, originalURL, err)
	//if err != nil {
	//	http.Error(w, "недопустимый идентификатор URL-адреса", http.StatusNotFound)
	//	return
	//}
	//fmt.Fprintf(w, "<script>location='%s';</script>", originalURL)

	id := r.URL.Path
	log.Printf("URL.Path: %s", id)
	id = id[1:]
	//log.Printf("URL.Path: %s", id)
	//if r.Method != http.MethodGet {
	//	http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
	//	return
	//}
	////id := chi.URLParam(r, "id")
	////log.Printf("id- %s", id)
	////if id == "" {
	////	http.Error(w, "параметр id пуст", http.StatusBadRequest)
	////	return
	////}
	//
	////Добавляем схему протокола, если она отсутствует
	//if !strings.HasPrefix(id, "http://") && !strings.HasPrefix(id, "https://") {
	//	id = "http://" + id
	//}
	//
	//parsedURL, err := url.Parse(id)
	//if err != nil {
	//	http.Error(w, "недопустимый формат URL-адреса", http.StatusBadRequest)
	//	return
	//}

	//Декодируем URL
	//decodedURL := parsedURL.String()
	//decodedURL, err = url.PathUnescape(decodedURL)
	//if err != nil {
	//	http.Error(w, "ошибка декодирования URL-адреса", http.StatusBadRequest)
	//	return
	//}
	//Ищем оригинальный URL в urlMap
	/*originalURL, ok := maintenance.NewMap().Get(id)
	//log.Printf("Извлечен URL из urlMap. Ключ: %s, Значение: %s, Найден: %v", decodedURL, originalURL, ok)
	if ok != nil {
		http.Error(w, "недопустимый идентификатор URL-адреса", http.StatusNotFound)
		return
	}*/
	originalURL, ok := urlmap[id]
	log.Printf("originalURL: %s", originalURL)
	if !ok {
		http.Error(w, "url не найден", http.StatusBadRequest)
	}
	// Возвращаем оригинальный URL
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
