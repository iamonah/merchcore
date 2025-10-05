package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

type tokenscope string

const (
	ActivationToken tokenscope = "activationToken"
	PasswordReset   tokenscope = "passwordReset"
)

type Token struct {
	Plaintext string
	TokenHash []byte
	UserID    uuid.UUID
	Expiry    time.Time
	Scope     tokenscope
}

func GenerateOTP(userId uuid.UUID, expiry time.Duration, scope tokenscope) (*Token, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}
	otp := fmt.Sprintf("%06d", n)

	hashOtp := sha256.Sum256([]byte(otp))
	return &Token{
		Plaintext: otp,
		TokenHash: hashOtp[:],
		UserID:    userId,
		Expiry:    time.Now().Add(expiry),
		Scope:     scope,
	}, nil
}

func GenerateToken(userID uuid.UUID, ttl time.Duration, scope tokenscope) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.TokenHash = hash[:]

	return token, nil
}
