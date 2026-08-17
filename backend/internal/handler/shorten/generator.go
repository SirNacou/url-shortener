package shorten

import (
	"crypto/rand"
	"fmt"
)

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate code")
	}

	for i, b := range bytes {
		bytes[i] = base62Alphabet[int(b)%len(base62Alphabet)]
	}

	return string(bytes), nil
}
