package service

import (
	"errors"
	"testing"
)

type mockGenerator struct{}

func (g *mockGenerator) GenerateRandom() string {
	return "mockkey"
}

type mockURLHandler struct {
	urlMap map[string]string
}

func (h *mockURLHandler) Add(keyURL, urlString string) {
	h.urlMap[keyURL] = urlString
}

func (h *mockURLHandler) Get(keyURL string) (string, error) {
	if urlString, ok := h.urlMap[keyURL]; ok {
		return urlString, nil
	}
	return "", errors.New("URL не найден")
}

func TestService_ShortenedURL(t *testing.T) {
	generator := &mockGenerator{}
	urlHandler := &mockURLHandler{urlMap: make(map[string]string)}
	service := NewService(generator, urlHandler)

	// Тестирование с непустым URL-адресом
	url := "https://example.com"
	expectedKey := "mockkey"
	expectedURL := "https://example.com"
	key, err := service.ShortenedURL(url)
	if err != nil {
		t.Errorf("Ошибка при сокращении URL: %v", err)
	}
	if key != expectedKey {
		t.Errorf("Ожидаемый ключ: %s, полученный ключ: %s", expectedKey, key)
	}
	if urlHandler.urlMap[key] != expectedURL {
		t.Errorf("Ожидаемый URL: %s, полученный URL: %s", expectedURL, urlHandler.urlMap[key])
	}

	// Тестирование с пустым URL-адресом
	emptyURL := ""
	_, err = service.ShortenedURL(emptyURL)
	if err == nil {
		t.Error("Ожидалась ошибка при передаче пустого URL-адреса")
	}
}

func TestService_IncreaseURL(t *testing.T) {
	urlHandler := &mockURLHandler{
		urlMap: map[string]string{
			"mockkey": "https://example.com",
		},
	}
	service := NewService(nil, urlHandler)

	// Тестирование с существующим ключом URL
	key := "mockkey"
	expectedURL := "https://example.com"
	url, err := service.IncreaseURL(key)
	if err != nil {
		t.Errorf("Ошибка при получении URL: %v", err)
	}
	if url != expectedURL {
		t.Errorf("Ожидаемый URL: %s, полученный URL: %s", expectedURL, url)
	}

	// Тестирование с несуществующим ключом URL
	nonexistentKey := "nonexistentkey"
	_, err = service.IncreaseURL(nonexistentKey)
	if err == nil {
		t.Error("Ожидалась ошибка при получении несуществующего ключа URL")
	}
}
