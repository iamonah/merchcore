package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrPhoneNumberExists   = errors.New("phone number already exists")
	ErrUserIDConflict      = errors.New("user with this ID already exists")
	ErrProviderIDExists    = errors.New("provider ID already exists")
	ErrProviderFieldsCheck = errors.New("invalid fields for provider type")
	ErrDatabase            = errors.New("database error")
	ErrTokenNotFound       = errors.New("otp not found")
	ErrTokenExpired        = errors.New("otp is expired")
	ErrSessionNotFound     = errors.New("user session not found")
)

type UserRepository interface {
	CreateUser(context.Context, *User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserPhoneNumber(ctx context.Context, phoneNum string) error
	UpdateUser(ctx context.Context, user *User) error
	VerifyUser(ctx context.Context, userID uuid.UUID) error
	CreateSession(ctx context.Context, s *Session) error
	GetSession(ctx context.Context, sessionId string) (*Session, error)
	CreateToken(ctx context.Context, otp *Token) error
	GetUserIDByToken(ctx context.Context, hash []byte, scope string) (uuid.UUID, error)
	DeleteToken(ctx context.Context, hash []byte, scope string) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error
	BlockSession(ctx context.Context, token []byte) error
}
