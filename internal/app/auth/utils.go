package auth

import (
	"time"

	"golang.org/x/exp/rand"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateAlphanumericCode(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
