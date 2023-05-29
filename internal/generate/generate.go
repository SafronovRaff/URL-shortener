package generate

import (
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    //  количество бит, которые будут использоваться для представления индекса символа
	letterIdxMask = 1<<letterIdxBits - 1 //  маска, используемая для извлечения битов индекса символа
	letterIdxMax  = 63 / letterIdxBits   // определяет максимальное количество индексов символов, которое помещается в 63 бита
)

type generate struct {
}

func NewGenerate() *generate {
	return &generate{}
}

func (g *generate) GenerateRandom() string {
	return GenerateRandomString(10)
}

var src = rand.NewSource(time.Now().UnixNano()) //используется для создания генератора случайных чисел на основе текущего времени

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
