package infrastructure

import (
	"github.com/alexedwards/argon2id"
)

type Argon2PasswordsHasher struct {
	Params *argon2id.Params
}

var argon2Instance *Argon2PasswordsHasher

func GetArgon2PasswordsHasher() *Argon2PasswordsHasher {
	if argon2Instance == nil {
		argon2Instance = &Argon2PasswordsHasher{
			// Docs: https://www.rfc-editor.org/rfc/rfc9106.html#name-parameter-choice
			Params: &argon2id.Params{
				Memory:      128 * 1024,
				Iterations:  4,
				Parallelism: 4,
				SaltLength:  16,
				KeyLength:   32,
			},
		}
	}

	return argon2Instance
}

func (hasher *Argon2PasswordsHasher) HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, hasher.Params)
	return hash, err
}

func (hasher *Argon2PasswordsHasher) ComparePasswords(password string, hashedPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	return match, err
}
