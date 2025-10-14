package auth

import (
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/users"

	"github.com/google/uuid"
)

type UserCreateReq struct {
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Country     string `json:"country" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
}

type UserCreateResp struct {
	Message string `json:"message"`
}

type UserUpdateReq struct {
	Password    *string `json:"password"`
	PhoneNumber *string `json:"phone_number"`
}

type UserSignInResp struct {
	TokenType             string    `json:"token_type"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserResp  `json:"user"`
}

func toCreateUser(ur UserCreateReq) users.UserCreate {
	return users.UserCreate{
		Password:    ur.Password,
		Email:       ur.Email,
		PhoneNumber: ur.PhoneNumber,
		Country:     ur.Country,
		FirstName:   ur.FirstName,
		LastName:    ur.LastName,
	}
}

type UserSignInReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenReq struct {
	Token string `json:"token" validate:"required"`
}

type UserResp struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	Country        string    `json:"country"`
	Roles          string    `json:"roles"`
	IsVerified     bool      `json:"is_verified"`
	Provider       *string   `json:"provider,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsStoreCreated bool      `json:"is_store_created"`
	NumOfStore     int       `json:"num_of_store"`
}

func toUserResp(u users.User) UserResp {
	return UserResp{
		UserID:         u.UserID,
		Email:          u.Email.Address,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Phone:          u.Contact.Number,
		Country:        u.Contact.Country,
		Roles:          u.GetRole(),
		IsVerified:     u.IsVerified,
		Provider:       u.Provider,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAT,
		IsStoreCreated: u.IsStoreCreated,
		NumOfStore:     u.NumOfStore,
	}
}

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

type ResetPasswordReq struct {
	Token           string `json:"token"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ChangePasswordReq struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenResp struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
