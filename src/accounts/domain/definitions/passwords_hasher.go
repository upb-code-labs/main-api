package definitions

type PasswordsHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswords(password string, hashedPassword string) (bool, error)
}
