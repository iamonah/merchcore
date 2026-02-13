package users

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/sdk/authz"
	"github.com/iamonah/merchcore/internal/sdk/errs"
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
		return nil, err
	}
	otp := fmt.Sprintf("%06d", n.Int64())

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
		return nil, fmt.Errorf("generatetoken: %w", err)
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.TokenHash = hash[:]

	return token, nil
}

type tokenData struct {
	AccessToken    string
	AcessExpiresAt time.Time
	UserId         uuid.UUID
}

func (s *UserBusiness) RenewAccessToken(ctx context.Context, payload *authz.Payload) (tokenData, error) {
	session, err := s.storer.GetSession(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return tokenData{}, errs.NewDomainError(errs.Unauthenticated, errors.New("invalid or expired refresh token"))
		}
		return tokenData{}, fmt.Errorf("getsession: %w", err)
	}

	if session.IsBlocked {
		return tokenData{}, errs.NewDomainError(errs.Unauthenticated, errors.New("session is blocked"))
	}

	if payload.UserID != session.UserID {
		return tokenData{}, errs.NewDomainError(errs.Unauthenticated, errors.New("session user mismatch"))
	}

	jwtData := authz.NewJWTData(session.UserID, payload.RoleID, 15*time.Minute, s.config.Observability.ServiceName)
	accessToken, accessPayload, err := s.authz.GenerateToken(jwtData)
	if err != nil {
		return tokenData{}, fmt.Errorf("generatetoken: %w", err)
	}

	data := tokenData{
		AccessToken:    accessToken,
		AcessExpiresAt: accessPayload.ExpiresAt.Time,
		UserId:         session.UserID,
	}
	return data, nil
}

func (s *UserBusiness) BlockSession(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	shaToken := sha256.Sum256([]byte(refreshToken))

	acckey := SessionAccessKeys(userID)
	err := s.cache.Delete(ctx, acckey)
	if err != nil {
		return fmt.Errorf("deletesession: %w", err)
	}

	refkey := SessionAccessKeys(userID)
	err = s.cache.Delete(ctx, refkey)
	if err != nil {
		return fmt.Errorf("deletesession: %w", err)
	}

	err = s.storer.BlockSession(ctx, shaToken[:])
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return errs.NewDomainError(errs.Unauthenticated, err)
		}
		return fmt.Errorf("blocksession: %w", err)
	}
	return nil
}
