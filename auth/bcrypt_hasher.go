package auth

import (
	"code.google.com/p/go.crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int
}

func (b BcryptHasher) Hash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), b.Cost)
	return string(hash)
}

func (b BcryptHasher) Check(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{Cost: bcrypt.DefaultCost}
}
