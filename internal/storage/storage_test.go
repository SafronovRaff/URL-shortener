package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_urlMap_Add(t *testing.T) {
	storage := NewStorage()

	t.Run("добавить URL", func(t *testing.T) {
		keyURL := "kjdnakdka"
		urlString := "http://example.com"
		storage.Add(keyURL, urlString)
		// проверка наличия добавленного URL
		resultURL, err := storage.Get(keyURL)
		assert.NoError(t, err, "неожиданная ошибка")
		assert.Equal(t, urlString, resultURL, "неожиданный URl")
	})
}

func Test_urlMap_Get(t *testing.T) {
	storage := NewStorage()

	t.Run("Получить существующий URL", func(t *testing.T) {
		keyURL := "kjdnakdka"
		urlString := "http://example.com"
		storage.Add(keyURL, urlString)
		// проверка получения существующего URL
		resultURL, err := storage.Get(keyURL)
		assert.NoError(t, err, "неожиданная ошибка")
		assert.Equal(t, urlString, resultURL, "неожиданный URl")
	})
	t.Run("Получить несуществующий URL", func(t *testing.T) {
		//проверка получения несуществующего URL
		_, err := storage.Get("nonexistent")
		assert.Error(t, err, "неожиданная ошибка")
		assert.Equal(t, "url не найден", err.Error(), "неожиданное сообщение об ошибке")
	})
}
