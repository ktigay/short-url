package generator

import (
	"github.com/ktigay/short-url/internal/random"
)

type RandStringGenerator struct {
}

func NewRandStringGenerator() *RandStringGenerator {
	return &RandStringGenerator{}
}

// Generate - рандомная строка длины n
func (s *RandStringGenerator) Generate(min, max int) string {
	return random.RandString(min, max)
}
