package storage

import (
	"errors"
	"log"
	"sync"
)

type urlMap struct {
	urlmap map[string]string //используется для хранения сокращенных URL на исходных URL
	mu     sync.Mutex
}

func NewMap() *urlMap {
	return &urlMap{
		urlmap: make(map[string]string),
	}
}

func (u *urlMap) Add(keyURL, urlString string) string {
	u.mu.Lock()
	u.urlmap[keyURL] = urlString
	log.Printf("Добавлен URL в urlMap. Ключ: %s, Значение: %s", keyURL, urlString)
	u.mu.Unlock()
	return urlString
}

func (u *urlMap) Get(keyURL string) (string, error) {
	u.mu.Lock()
	orl, ok := u.urlmap[keyURL]
	if !ok {
		return "", errors.New("url не найден")
	}
	u.mu.Unlock()
	return orl, nil
}
