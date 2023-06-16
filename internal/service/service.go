package service

import (
	"log"
	"net/url"
)

type generate interface {
	GenerateRandom() string
}

type URLhand interface {
	Add(keyURL, urlString string)
	Get(keyURL string) (string, error)
}

type Service struct {
	generate   generate
	urlHandler URLhand
}

func NewService(generate generate, urlHandler URLhand) *Service {
	return &Service{
		generate:   generate,
		urlHandler: urlHandler,
	}
}

func (s *Service) ShortenedURL(b string) (string, error) {
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
	s.urlHandler.Add(keyURL, urlString)
	return keyURL, err
}

func (s *Service) IncreaseURL(id string) (string, error) {
	originalURL, err := s.urlHandler.Get(id)
	if err != nil {
		log.Printf("НЕ найден URL в urlMap. Ключ: %s", originalURL)
		return "", err
	}
	return originalURL, nil
}
