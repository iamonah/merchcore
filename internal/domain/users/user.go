package users

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/shared/contact"
	"github.com/IamOnah/storefronthq/internal/domain/shared/role"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"

	"github.com/google/uuid"
)

type Provider string

const Local Provider = "local"
const Google Provider = "google"

type User struct {
	UserID         uuid.UUID
	PasswordHash   []byte
	Email          *mail.Address
	FirstName      string
	LastName       string
	Contact        contact.Contact
	CreatedAt      time.Time
	UpdatedAT      time.Time
	Provider       *string //(local or google)
	ProviderID     *string
	Role           role.Role
	IsVerified     bool
	DeletedAt      *time.Time
	IsStoreCreated bool
	IsEnabled      bool
	NumOfStore     int
}

func generateUserID() uuid.UUID {
	return uuid.New()
}

func NewUser(userInfo UserCreate) (User, error) {
	fieldErrs := errs.NewFieldErrors()

	user := User{
		UserID:     generateUserID(),
		IsVerified: false,
		IsEnabled:  false,
	}

	cleanFirstName := strings.TrimSpace(userInfo.FirstName)
	if cleanFirstName == "" {
		fieldErrs.AddFieldError("first_name", errors.New("cannot be empty"))
	}

	cleanLastName := strings.TrimSpace(userInfo.LastName)
	if cleanLastName == "" {
		fieldErrs.AddFieldError("last_name", errors.New("cannot be empty"))
	}

	newContact := contact.NewContact(userInfo.PhoneNumber, userInfo.Country)
	if err := newContact.ValidateContact(); err != nil {
		fieldErrs.AddFieldError("phone_number", err)
	}

	email, err := mail.ParseAddress(userInfo.Email)
	if err != nil {
		fieldErrs.AddFieldError("email", err)
	}

	password, err := HashPassword([]byte(userInfo.Password))
	if err != nil {
		fieldErrs.AddFieldError("password", err)
	}
	user.PasswordHash = password

	if err := fieldErrs.ToError(); err != nil {
		return User{}, err
	}

	// provider logic
	if userInfo.ProviderID != nil {
		user.Provider = userInfo.Provider
		user.ProviderID = userInfo.ProviderID
	} else {
		provider := string(Local)
		user.Provider = &provider
	}

	user.Email = email
	user.Contact = newContact
	user.LastName = cleanLastName
	user.FirstName = cleanFirstName
	user.CreatedAt = time.Now()
	user.UpdatedAT = time.Now()
	user.NumOfStore = 0
	user.IsStoreCreated = false
	user.Role = role.Guest

	return user, nil
}

func (u *User) GetPhoneNumber() string {
	return u.Contact.Number
}

func (u *User) GetEmail() string {
	return u.Email.Address
}

func (u *User) GetRole() string {
	return u.Role.String()
}
