package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// DefaultCost is recommended
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	// Compares the plaintext password with the stored hash.
	// Returns nil on success, or an error on failure.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
