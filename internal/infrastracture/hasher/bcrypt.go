package hasher

import "golang.org/x/crypto/bcrypt"

func BcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func BcryptComparePasswordAndHash(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func BcryptIsInvalidPasswordError(err error) bool {
	return err == bcrypt.ErrMismatchedHashAndPassword
}
