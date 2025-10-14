package users

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const HashCost = 12

var ErrInvalidPassword = errors.New("incorrect password")

func HashPassword(password []byte) ([]byte, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), HashCost)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrHashTooShort):
			return nil, fmt.Errorf("hash too short: %w", err)
		}
		return nil, err
	}
	return hashPassword, nil
}

func ComparePassword(storedHash, password []byte) error {
	err := bcrypt.CompareHashAndPassword(storedHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}

