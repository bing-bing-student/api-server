package utils

import "crypto/rand"

func GenerateRandomDigits(length int) (string, error) {
	const digits = "0123456789"
	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	for i, b := range buffer {
		buffer[i] = digits[b%10]
	}

	return string(buffer), nil
}
