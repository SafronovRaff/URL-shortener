package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockServices struct{}

func (m *mockServices) IncreaseURL(string) (string, error) {
	return "http://example.com", nil
}

func (m *mockServices) ShortenedURL(string) (string, error) {
	return "kjdnakdka", nil
}

func TestHandlers_Shortened(t *testing.T) {
	// handler с фиктивным сервисом
	h := NewHandlers(&mockServices{})
	t.Run("Действительный запрос", func(t *testing.T) {
		//создаём фиктивный запрос для записи в тело
		body := "http://example.com/"
		req := httptest.NewRequest(http.MethodPost, "/shortened", bytes.NewBufferString(body))

		// переменная для запись фиктивного ответа на сервер
		recorder := httptest.NewRecorder()
		// обработчик
		h.Shortened(recorder, req)
		// проверяем код состояния ответа
		assert.Equal(t, http.StatusCreated, recorder.Code, "неожиданный код состояния, ожидал 201 ")
		// проверяем заголовок
		assert.Equal(t, "text/plain; charset=utf-8", recorder.Header().Get("Content-Type"), "неожиданный Content-Type")
		expectedResponse := "http://localhost:8080/kjdnakdka"
		assert.Equal(t, expectedResponse, recorder.Body.String(), "неожиданное тело ответа")
	})

}

func TestHandlers_Increase(t *testing.T) {
	// handler с фиктивным сервисом
	h := NewHandlers(&mockServices{})

	t.Run("действительный запрос", func(t *testing.T) {
		//фиктивный Get с добавлением переменно	 пути
		req := httptest.NewRequest(http.MethodGet, "/increase/kjdnakdka", nil)

		recorder := httptest.NewRecorder()

		// обработчик
		h.Increase(recorder, req)
		// проверяем код состояния ответа
		assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code, "неожиданный код состояния, ожидал 307")
		// проверяем заголовок Location
		assert.Equal(t, "http://example.com", recorder.Header().Get("Location"), "неожиданная локация заголовка")
	})

	t.Run("Неверный метод", func(t *testing.T) {
		// cоздаем тестовый запрос с недопустимым методом
		req := httptest.NewRequest(http.MethodPost, "/increase/kjdnakdka", nil)
		recorder := httptest.NewRecorder()
		// обработчик
		h.Increase(recorder, req)

		// gроверка кода состояния ответа
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "unexpected status code")
	})

}
func TestHandlers_ShortenHandler(t *testing.T) {
	// handler с фиктивным сервисом
	h := NewHandlers(&mockServices{})

	t.Run("действительный запрос", func(t *testing.T) {
		// Создаем тестовый JSON-запрос
		requestBody := `{"URL":"http://example.com"}`
		// Создаем тестовый запрос с методом POST и правильным Content-Type
		req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		// Обработчик
		h.ShortenHandler(recorder, req)
		// Проверяем код состояния ответа
		assert.Equal(t, http.StatusOK, recorder.Code, "неожиданный код состояния, ожидал 200")
		// Проверяем Content-Type заголовок
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "неожиданный Content-Type заголовка")
		// Читаем тело ответа
		var resp ShortenResponse
		err := json.NewDecoder(recorder.Body).Decode(&resp)
		assert.NoError(t, err, "ошибка декодирования JSON-ответа")
		// Проверяем, что декодирование JSON-ответа прошло успешно
		assert.Equal(t, "kjdnakdka", resp.Result, "неожиданный результат")

	})

	t.Run("недопустимый метод запроса", func(t *testing.T) {
		// Создаем тестовый запрос с методом GET
		req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
		recorder := httptest.NewRecorder()
		// Обработчик
		h.ShortenHandler(recorder, req)
		// Проверяем код состояния ответа
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "неожиданный код состояния, ожидал 400")
	})

	t.Run("недопустимый Content-Type", func(t *testing.T) {
		// Создаем тестовый запрос с методом POST, но неправильным Content-Type
		req := httptest.NewRequest(http.MethodPost, "/shorten", nil)
		req.Header.Set("Content-Type", "text/plain")
		recorder := httptest.NewRecorder()
		// Обработчик
		h.ShortenHandler(recorder, req)
		// Проверяем код состояния ответа
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "неожиданный код состояния, ожидал 400")
	})

	t.Run("не верное тело запроса", func(t *testing.T) {
		// Создаем тестовый JSON-запрос с неверным форматом
		requestBody := `{"invalidField":"http://example.com"}`
		// Создаем тестовый запрос с методом POST и правильным Content-Type
		req := httptest.NewRequest(http.MethodPut, "/shorten", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		// Обработчик
		h.ShortenHandler(recorder, req)
		// Проверяем код состояния ответа
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "неожиданный код состояния, ожидал 400")
	})

}
