package users

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const HashCost = 12

func HashPassword(password []byte) ([]byte, error) {
	sha := sha256.Sum256(password)
	shaHex := hex.EncodeToString(sha[:]) //eliminate the password-to-long error

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(shaHex), HashCost)
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
	sha := sha256.Sum256(storedHash)
	shaHex := hex.EncodeToString(sha[:])

	err := bcrypt.CompareHashAndPassword([]byte(shaHex), password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("incorrect password: %w", err)
		}
		return err
	}
	return nil
}
