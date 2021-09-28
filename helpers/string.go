package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomToken(length int) (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
