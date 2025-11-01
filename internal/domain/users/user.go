package users

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/iamonah/merchcore/internal/domain/shared/contact"
	"github.com/iamonah/merchcore/internal/domain/shared/role"
	"github.com/iamonah/merchcore/internal/sdk/errs"

	"github.com/google/uuid"
)

type Provider string

var (
	Local  = newProvider("local")
	Google = newProvider("google")
)

func (p Provider) String() string {
	return string(p)
}

var providers = make(map[string]Provider)

func newProvider(v string) Provider {
	provider := Provider(v)
	providers[v] = provider
	return provider
}

func ParseProvider(v string) (Provider, error) {
	provider, ok := providers[v]
	if !ok {
		return "", fmt.Errorf("invalid provider: %s", v)
	}
	return provider, nil
}

type User struct {
	UserID         uuid.UUID
	PasswordHash   []byte
	Email          *mail.Address
	FirstName      string
	LastName       string
	Contact        contact.Contact
	CreatedAt      time.Time
	UpdatedAT      time.Time
	Provider       Provider //(local or google)
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
	userInfo.Sanitize()
	fieldErrs := errs.NewFieldErrors()

	user := User{
		UserID:         generateUserID(),
		IsVerified:     false,
		IsEnabled:      false,
		CreatedAt:      time.Now(),
		UpdatedAT:      time.Now(),
		NumOfStore:     0,
		IsStoreCreated: false,
		Role:           role.Guest,
	}

	if userInfo.FirstName == "" {
		fieldErrs.AddFieldError("first_name", errors.New("cannot be empty"))
	}

	if userInfo.LastName == "" {
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

	if userInfo.Password == "" {
		fieldErrs.AddFieldError("password", errors.New("cannot be empty"))
	} else {
		hashed, err := HashPassword([]byte(userInfo.Password))
		if err != nil {
			return User{}, fmt.Errorf("hashpassword: %w", err)
		}
		user.PasswordHash = hashed
	}

	if err := fieldErrs.ToError(); err != nil {
		return User{}, err
	}

	if userInfo.ProviderID != nil && userInfo.Provider != nil {
		provider, err := ParseProvider(*userInfo.Provider)
		if err != nil {
			return User{}, fmt.Errorf("parseprovider: %w", err)
		}
		user.Provider = provider
		user.ProviderID = userInfo.ProviderID
	} else {
		//used password and email signup
		user.Provider = Local
	}

	user.Email = email
	user.Contact = newContact
	user.LastName = userInfo.LastName
	user.FirstName = userInfo.FirstName

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
