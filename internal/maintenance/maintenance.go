package maintenance

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    //  количество бит, которые будут использоваться для представления индекса символа
	letterIdxMask = 1<<letterIdxBits - 1 //  маска, используемая для извлечения битов индекса символа
	letterIdxMax  = 63 / letterIdxBits   // определяет максимальное количество индексов символов, которое помещается в 63 бита
)

var src = rand.NewSource(time.Now().UnixNano()) //используется для создания генератора случайных чисел на основе текущего времени

type urlMap struct {
	urlmap map[string]string //используется для хранения сокращенных URL на исходных URL
	mu     sync.Mutex
}

// Функция генерирует случайную строку длиной "n" из  байтового слайса
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func NewMap() *urlMap {
	return &urlMap{urlmap: make(map[string]string)}
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
