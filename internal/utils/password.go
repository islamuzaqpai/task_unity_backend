package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//затестить

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("failed to hash a password: %w", err)
	}

	return string(hashedPassword), nil
}

func ValidatePassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, fmt.Errorf("failed to validate passwords: %w", err)
	}

	return true, nil
}
