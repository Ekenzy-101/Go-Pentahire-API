package helpers

import (
	"crypto/rand"
	"math/big"
)

const (
	numbers = "0123456789"
	letters = numbers + "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

func GenerateRandomNumbers(length int) (string, error) {
	return generateRandomFromSeed(length, numbers)
}

func GenerateRandomToken(length int) (string, error) {
	return generateRandomFromSeed(length, letters)
}

func generateRandomFromSeed(length int, seed string) (string, error) {
	byteSlice := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(seed))))
		if err != nil {
			return "", err
		}
		byteSlice[i] = seed[num.Int64()]
	}
	return string(byteSlice), nil
}
