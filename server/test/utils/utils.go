package utils

import (
	"math/rand/v2"
)

func MakeRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.IntN(len(letters))]
	}

	return string(result)
}
