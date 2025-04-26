package random

import (
	"math/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(min, max int) string {
	var n int
	if min == max {
		n = min
	} else {
		n = rand.Intn(max-min) + min
	}

	l := len(charset)
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(l)])
	}
	return sb.String()
}
