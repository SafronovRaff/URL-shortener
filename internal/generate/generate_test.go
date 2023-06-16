package generate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	generator := NewGenerate()
	t.Run("Создание рандомной строки", func(t *testing.T) {
		res := generator.GenerateRandom()
		//проверка длинны строки
		assert.Equal(t, 10, len(res), "неожиданная длина строки")
		//проверка символов в строке
		for _, s := range res {
			assert.Contains(t, letterBytes, string(s), "неожиданный символ")
		}

	})

}
