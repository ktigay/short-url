package generator

import (
	"math/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type RandStringGenerator struct {
}

func NewRandStringGenerator() *RandStringGenerator {
	return &RandStringGenerator{}
}

// Generate - рандомная строка длины n
func (s *RandStringGenerator) Generate(min int, max int) string {
	n := rand.Intn(max-min) + min
	l := len(charset)
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(l)])
	}
	return sb.String()
}
