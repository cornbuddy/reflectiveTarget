package entities

import (
	"github.com/alexedwards/argon2id"
)

type Password struct {
	Hash string
}

func (p Password) Verify(plaintext string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plaintext, p.Hash)
	if err != nil {
		return false, err
	}

	return match, nil
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
