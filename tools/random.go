package tools

import (
	"math/rand"
	"time"
)

func NewRandomString(size int, alph string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune(alph)

	b := make([]rune, size)

	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
