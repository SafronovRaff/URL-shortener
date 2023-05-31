package storage

import (
	"errors"
	"log"
	"sync"
	"unsafe"
)

func NewStorage() *urlMap {
	return &urlMap{
		urlmap: make(map[string]string),
	}
}

type urlMap struct {
	urlmap map[string]string //используется для хранения сокращенных URL на исходных URL
	mu     sync.RWMutex
}

func (u *urlMap) Add(keyURL, urlString string) {
	u.mu.Lock()
	u.urlmap[keyURL] = urlString
	log.Printf("Добавлен URL в urlMap. Ключ: %s, Значение: %s", keyURL, urlString)
	log.Printf("Адрес мапы при Add`urlmap`: %p \n", unsafe.Pointer(&u.urlmap))
	u.mu.Unlock()
}

func (u *urlMap) Get(id string) (string, error) {
	u.mu.RLock()
	log.Printf("Адрес мапы при Get `urlmap`: %p \n", unsafe.Pointer(&u.urlmap))
	defer u.mu.RUnlock()
	originalURL, ok := u.urlmap[id]
	if !ok {
		return "", errors.New("url не найден")
	}

	return originalURL, nil
}
