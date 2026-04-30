package utils

import (
	"encoding/json"
	"math/rand/v2"
)

func DeepCopy[T any](val T) (*T, error) {
	bytes, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	var res T
	json.Unmarshal(bytes, &res)

	return &res, nil
}

func MakeRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.IntN(len(letters))]
	}

	return string(result)
}
