package users

import (
	"net/mail"
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/shared/contact"
	"github.com/IamOnah/storefronthq/internal/domain/shared/role"
	"github.com/google/uuid"
)

type UserCreate struct {
	Password    string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Country     string
	Provider    *string
	ProviderID  *string
}

type UserUpdate struct {
	Password    *string
	PhoneNumber *string
}

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken []byte
	UserAgent    string
	ClientIP     string
	IsBlocked    bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func (ses *Session) IsSessionBlocked() bool {
	return ses.IsBlocked
}

func (ses *Session) UserIDCheck(UserId uuid.UUID) bool {
	return ses.UserID != UserId
}

func (ses *Session) IsSessionExpired() bool {
	return time.Now().After(ses.ExpiresAt)
}

type UpdateUser struct {
	FirstName      *string
	LastName       *string
	Email          *mail.Address
	Roles          *role.Role
	Password       *string
	IsEnabled      *bool
	PhoneNumber    *contact.Contact
	UpdatedAt      *time.Time
	NumOfStore     *int
	IsStoreCreated *bool
}
