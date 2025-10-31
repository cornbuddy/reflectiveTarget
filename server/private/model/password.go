package model

import (
	"github.com/alexedwards/argon2id"
)

type Password struct {
	Hash string
}

func NewPassword(plaintext string) (*Password, error) {
	params := &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  64,
		Parallelism: 4,
		SaltLength:  32,
		KeyLength:   32,
	}
	hash, err := argon2id.CreateHash(plaintext, params)
	if err != nil {
		return nil, err
	}

	return &Password{Hash: hash}, nil
}
