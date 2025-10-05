package users

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/shared/phone"
	"github.com/IamOnah/storefronthq/internal/domain/shared/role"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"

	"github.com/google/uuid"
)

// some work
type User struct {
	UserID         uuid.UUID
	PasswordHash   []byte
	Email          *mail.Address
	FirstName      string
	LastName       string
	PhoneNumber    phone.PhoneNumber
	CreatedAt      time.Time
	UpdatedAT      time.Time
	Provider       *string //(local or google)
	ProviderID     *string
	Roles          role.Role
	IsVerified     bool
	Country        string
	DeletedAt      *time.Time
	IsStoreCreated bool
}

func generateUserID() uuid.UUID {
	return uuid.New()
}

func NewUser(userInfo UserCreate) (User, error) {
	fieldErrs := errs.NewFieldErrors()

	user := User{
		UserID:     generateUserID(),
		IsVerified: false,
	}

	cleanFirstName := strings.TrimSpace(userInfo.FirstName)
	if cleanFirstName == "" {
		fieldErrs.AddFieldError("last_name", errors.New("cannot be emtpy"))
	}

	cleanLastName := strings.TrimSpace(userInfo.LastName)
	if cleanLastName == "" {
		fieldErrs.AddFieldError("fist_name", errors.New("cannot be empty"))
	}

	newPhoneNum := phone.NewPhoneNumber(userInfo.PhoneNumber, userInfo.Country)
	err := newPhoneNum.ValidateNumber()
	if err != nil {
		fieldErrs.AddFieldError("phone_number", fmt.Errorf("invalid: %w", err))
	}

	email, err := mail.ParseAddress(userInfo.Email)
	if err != nil {
		fieldErrs.AddFieldError("email", fmt.Errorf("invalid: %w", err))
	}

	if userInfo.Password != nil {
		password, err := authz.HashPassword([]byte(*userInfo.Password))
		if err != nil {
			fieldErrs.AddFieldError("password", fmt.Errorf("invalid: %w", err))
		}
		user.PasswordHash = password
		*user.Provider = "local"
	}

	if err := fieldErrs.ToError(); err != nil {
		return User{}, err
	}
	//oauth0 provider id and provider set
	if userInfo.ProviderID != nil {
		user.Provider = userInfo.Provider
		user.ProviderID = userInfo.ProviderID
	}

	user.Email = email
	user.PhoneNumber = newPhoneNum
	user.LastName = cleanLastName
	user.FirstName = cleanFirstName
	user.CreatedAt = time.Now()
	user.UpdatedAT = time.Now()

	return user, nil
}

func (u *User) GetPhoneNumber() string {
	return u.PhoneNumber.Number
}

func (u *User) GetEmail() string {
	return u.Email.Address
}

func (u *User) GetRole() string {
	return u.Roles.Name
}

func (u *User) GetUserPermissions() []role.Permission {
	return u.Roles.Permissions
}
