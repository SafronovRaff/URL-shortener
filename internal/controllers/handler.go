package controllers

import (
	"github.com/SafronovRaff/URL-shortener/internal/maintenance"
	"io"
	"log"
	"net/http"
	"net/url"
)

var urlmap map[string]string

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
	log.Printf("найден originalURL: %s", originalURL)
	if !ok {
		http.Error(w, "originalURL не найден", http.StatusBadRequest)
	}
	// Возвращаем оригинальный URL

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
