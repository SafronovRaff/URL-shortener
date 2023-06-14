package controllers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
		/*	//задаем значение переменной "id" в URL
			req = mux.SetURLVars(req, map[string]string{"id": "kjdnakdka"})
			assert.NoError(t, err, "не удалось создать запрос")
			// переменная для запись фиктивного ответа на сервер
			recorder := httptest.NewRecorder()
		*/

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
