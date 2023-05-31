package service

import (
	"log"
	"net/url"
)

func (s *service) ShortenedURL(b string) (string, error) {
	// Преобразуем данные в строку URL
	urlString, err := url.PathUnescape(string(b))
	if err != nil {
		return "", err
	}
	// Вывод значения URL в лог
	log.Printf("Извлеченное значение URL: %s", urlString)
	// Генерация случайной строки в качестве ключа
	keyURL := s.generate.GenerateRandom()
	// Добавление значения URL в urlMap
	s.urlhandler.Add(keyURL, urlString)
	return keyURL, err
}

func (s *service) IncreaseURL(id string) (string, error) {
	originalURL, err := s.urlhandler.Get(id)
	if err != nil {
		log.Printf("НЕ найден URL в urlMap. Ключ: %s", originalURL)
		return "", err
	}
	return originalURL, nil
}

type generate interface {
	GenerateRandom() string
}

type service struct {
	generate   generate
	urlhandler URLhand
}

func NewService(generate generate, urlhandler URLhand) *service {
	return &service{
		generate:   generate,
		urlhandler: urlhandler,
	}
}

type URLhand interface {
	Add(keyURL, urlString string)
	Get(keyURL string) (string, error)
}
