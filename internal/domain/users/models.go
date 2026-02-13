package users

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/domain/types/contact"
	"github.com/iamonah/merchcore/internal/domain/types/role"
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

func (u *UserCreate) Sanitize() {
	u.Password = strings.TrimSpace(u.Password)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.PhoneNumber = strings.TrimSpace(u.PhoneNumber)
	u.Country = strings.TrimSpace(u.Country)

	if u.Provider != nil {
		provider := strings.TrimSpace(*u.Provider)
		u.Provider = &provider
	}

	if u.ProviderID != nil {
		providerID := strings.TrimSpace(*u.ProviderID)
		u.ProviderID = &providerID
	}
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
