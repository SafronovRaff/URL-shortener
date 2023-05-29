package storage

import (
	"errors"
	"log"
	"sync"
)

func NewStorage() *urlMap {
	return &urlMap{
		urlmap: make(map[string]string),
		mu:     sync.Mutex{},
	}
}

type urlMap struct {
	urlmap map[string]string //используется для хранения сокращенных URL на исходных URL
	mu     sync.Mutex
}

func (u *urlMap) Add(keyURL, urlString string) {
	u.mu.Lock()
	u.urlmap[keyURL] = urlString
	log.Printf("Добавлен URL в urlMap. Ключ: %s, Значение: %s", keyURL, urlString)
	u.mu.Unlock()
}

func (u *urlMap) Get(keyURL string) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	orl, ok := u.urlmap[keyURL]
	if !ok {

		return "", errors.New("url не найден")
	}

	return orl, nil
}
