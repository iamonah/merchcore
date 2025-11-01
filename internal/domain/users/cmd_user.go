package users

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/config"
	"github.com/iamonah/merchcore/internal/infra/cache"
	"github.com/iamonah/merchcore/internal/infra/database"
	"github.com/iamonah/merchcore/internal/sdk/authz"
	"github.com/iamonah/merchcore/internal/sdk/errs"
	// Global logger
)

type UserBusiness struct {
	storer UserRepository
	trx    database.TransactorTX
	authz  authz.TokenMaker
	config *config.Config
	cache  cache.Cache
}

func NewUserBusiness(store UserRepository, trx *database.TRXManager, authz *authz.JWTAuthMaker, cfg *config.Config) *UserBusiness {
	return &UserBusiness{
		storer: store,
		trx:    trx,
		authz:  authz,
		config: cfg,
	}
}

func (s *UserBusiness) CreateUser(ctx context.Context, info UserCreate) (User, Token, error) {
	user, err := NewUser(info)
	if err != nil {
		return User{}, Token{}, errs.NewDomainError(errs.InvalidArgument, err)
	}

	var token Token
	err = s.trx.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.storer.CreateUser(ctx, &user); err != nil {
			if errors.Is(err, ErrUserIDConflict) {
				// Retry with a new UUID once
				user.UserID = uuid.New()
				if err := s.storer.CreateUser(ctx, &user); err != nil {
					return fmt.Errorf("createuserretry: %w", err)
				}
			}
			return err
		}

		t, err := GenerateOTP(user.UserID, 90*time.Second, ActivationToken)
		if err != nil {
			return fmt.Errorf("generateOTP: %w", err)
		}

		if err := s.storer.CreateToken(ctx, t); err != nil {
			return fmt.Errorf("createOTP: %w", err)
		}
		token = *t
		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists), errors.Is(err, ErrPhoneNumberExists):
			return User{}, Token{}, errs.NewDomainError(errs.AlreadyExists, err)
		default:
			return User{}, Token{}, fmt.Errorf("createuser-trx: %w", err)
		}
	}

	return user, token, nil
}

func (s *UserBusiness) UpdateUser(ctx context.Context, userID uuid.UUID, uu UpdateUser) (*User, error) {
	usr, err := s.storer.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if uu.FirstName != nil {
		usr.FirstName = *uu.FirstName
	}

	if uu.Email != nil {
		usr.Email = uu.Email
	}

	if uu.Roles != nil {
		usr.Role = *uu.Roles
	}

	if uu.Password != nil {
		err := ComparePassword(usr.PasswordHash, []byte(*uu.Password))
		if err != nil {
			return nil, fmt.Errorf("comparepassword: %w", err)
		}
		pw, err := HashPassword([]byte(*uu.Password))
		if err != nil {
			return nil, fmt.Errorf("generate password hash: %w", err)
		}
		usr.PasswordHash = pw
	}

	if uu.IsEnabled != nil {
		usr.IsEnabled = *uu.IsEnabled
	}

	usr.UpdatedAT = time.Now()

	if err := s.storer.UpdateUser(ctx, usr); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return usr, nil
}

func (s *UserBusiness) ActivateUser(ctx context.Context, usrID uuid.UUID, token string) error {
	sha := sha256.Sum256([]byte(token))

	id, err := s.storer.GetUserIDByToken(ctx, sha[:], string(ActivationToken))
	if err != nil {
		if errors.Is(err, ErrDatabase) {
			return fmt.Errorf("getuseridbytoken: %w", err)
		}
		return errs.NewDomainError(errs.InvalidArgument, errors.New("invalid or expired token"))
	}
	if id != usrID {
		return errs.NewDomainError(errs.InvalidArgument, errors.New("invalid or expired token"))
	}

	if err := s.trx.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.storer.VerifyUser(ctx, id); err != nil {
			return err
		}
		if err := s.storer.DeleteToken(ctx, sha[:], string(ActivationToken)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, ErrDatabase) {
			return fmt.Errorf("activateuser-trx: %w", err)
		}
		return errs.NewDomainError(errs.InvalidArgument, errors.New("invalid or expired OTP"))
	}

	return nil
}

func (s *UserBusiness) ResendActivationToken(ctx context.Context, userID uuid.UUID) (User, Token, error) {
	u, err := s.storer.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, Token{}, errs.NewDomainError(errs.InvalidArgument, err)
		}
		return User{}, Token{}, fmt.Errorf("getuserbyid:%w", err)
	}

	if u.IsVerified {
		return User{}, Token{}, errs.NewDomainError(errs.InvalidArgument, errors.New("user already verified"))
	}

	t, err := GenerateOTP(userID, 90*time.Second, ActivationToken)
	if err != nil {
		return User{}, Token{}, fmt.Errorf("generateotp: %w", err)
	}
	if err := s.storer.CreateToken(ctx, t); err != nil {
		return User{}, Token{}, fmt.Errorf("createotp: %w", err)
	}

	return *u, *t, nil
}

func (s *UserBusiness) Authenticate(ctx context.Context, email *mail.Address, password string) (User, error) {
	user, err := s.storer.GetUserByEmail(ctx, email.Address)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, errs.NewDomainError(errs.Unauthenticated, errors.New("invalid credentials"))
		}
		return User{}, fmt.Errorf("getinguserbyemail: %w", err)
	}

	if len(user.PasswordHash) == 0 {
		return User{}, errs.NewDomainError(errs.Internal, errors.New("invalid user data"))
	}
	if err := ComparePassword(user.PasswordHash, []byte(password)); err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			return User{}, errs.NewDomainError(errs.Unauthenticated, errors.New("invalid credentials"))
		}
		return User{}, fmt.Errorf("comparepassword: %w", err)
	}
	return *user, nil
}

type SessionData struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (s *UserBusiness) CreateSession(ctx context.Context, user User, userAgent, clientIP string) (*SessionData, error) {
	accessTokenData := authz.NewJWTData(user.UserID, user.GetRole(), 30*time.Minute, s.config.Observability.ServiceName)
	accessToken, accessPayload, err := s.authz.GenerateToken(accessTokenData)
	if err != nil {
		return nil, fmt.Errorf("generatetoken: %w", err)
	}

	refreshTokenData := authz.NewJWTData(user.UserID, user.GetRole(), 24*time.Hour, s.config.Observability.ServiceName)
	refreshToken, refreshPayload, err := s.authz.GenerateToken(refreshTokenData)
	if err != nil {
		return nil, fmt.Errorf("generatetoken: %w", err)
	}
	//cache
	acckey := SessionAccessKeys(user.UserID)
	err = s.cache.Set(ctx, acckey, accessToken, time.Duration(accessPayload.ExpiresAt.Second()))
	if err != nil {
		return nil, fmt.Errorf("setcache: %w", err)
	}
	refkey := SessionRefreshKeys(user.UserID)
	err = s.cache.Set(ctx, refkey, accessToken, time.Duration(refreshPayload.ExpiresAt.Second()))
	if err != nil {
		return nil, fmt.Errorf("setcache: %w", err)
	}

	hashRefresh := sha256.Sum256([]byte(refreshToken))
	err = s.storer.CreateSession(ctx, &Session{
		ID:           uuid.MustParse(refreshPayload.ID),
		UserID:       user.UserID,
		RefreshToken: hashRefresh[:],
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		return nil, fmt.Errorf("createsession: %w", err)
	}

	data := &SessionData{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
	}
	return data, nil
}

func (s *UserBusiness) ForgetPassword(ctx context.Context, email *mail.Address) (*User, *Token, error) {
	user, err := s.storer.GetUserByEmail(ctx, email.Address)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Not an error: user doesn't exist
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("getuserbyemail: %w", err)
	}

	token, err := GenerateToken(user.UserID, 15*time.Minute, PasswordReset)
	if err != nil {
		return nil, nil, fmt.Errorf("generatetoken: %w", err)
	}

	fmt.Println(token.Plaintext)
	if err := s.storer.CreateToken(ctx, token); err != nil {
		return nil, nil, fmt.Errorf("createtoken: %w", err)
	}

	return user, token, nil
}

func (s *UserBusiness) PasswordReset(ctx context.Context, newPass string, token string) (uuid.UUID, error) {
	newPassword, err := HashPassword([]byte(newPass))
	if err != nil {
		return uuid.Nil, fmt.Errorf("hashpassword: %w", err)
	}

	shaToken := sha256.Sum256([]byte(token))
	var userID uuid.UUID

	err = s.trx.WithTransaction(ctx, func(ctx context.Context) error {
		id, err := s.storer.GetUserIDByToken(ctx, shaToken[:], string(PasswordReset))
		if err != nil {
			return fmt.Errorf("getuseridbytoken: %w", err)
		}
		userID = id

		if err := s.storer.DeleteToken(ctx, shaToken[:], string(PasswordReset)); err != nil {
			return fmt.Errorf("deletetoken: %w", err)
		}

		if err := s.storer.UpdatePassword(ctx, userID, newPassword); err != nil {
			return fmt.Errorf("updatepassword: %w", err)
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, ErrDatabase) {
			return uuid.Nil, fmt.Errorf("dbtransaction: %w", err)
		}
		return uuid.Nil, errs.NewDomainError(errs.Unauthenticated, errors.New("link expired or invalid"))
	}
	return userID, nil
}

func (s *UserBusiness) ChangePassword(ctx context.Context, userId uuid.UUID, oldPass, newPass string) (User, error) {
	var userData *User
	err := s.trx.WithTransaction(ctx, func(ctx context.Context) error {
		userValue, err := s.storer.GetUserByID(ctx, userId)
		if err != nil {
			return fmt.Errorf("getuserbyid: %w", err)
		}

		if err := ComparePassword(userValue.PasswordHash, []byte(oldPass)); err != nil {
			return fmt.Errorf("comparePassword: %w", err)
		}

		newHash, err := HashPassword([]byte(newPass))
		if err != nil {
			return fmt.Errorf("hashpassword: %w", err)
		}

		userValue.PasswordHash = newHash
		if err := s.storer.UpdateUser(ctx, userValue); err != nil {
			return err
		}

		userData = userValue
		return nil
	})

	// user not found db inconsistency cause user must be logged in to change password
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			return User{}, errs.NewDomainError(errs.Unauthenticated, errors.New("password incorrect"))
		}
		return User{}, fmt.Errorf("dbtransaction : %w", err)

	}
	return *userData, nil
}
