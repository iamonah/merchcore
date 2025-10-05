package userstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/IamOnah/storefronthq/internal/domain/users"
	"github.com/IamOnah/storefronthq/internal/infra/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	conn database.DBTX // default connection or pool
}

func NewUserStore(conn *pgxpool.Pool) *UserStore {
	return &UserStore{
		conn: conn,
	}
}
func (us *UserStore) CreateUser(ctx context.Context, usr *users.User) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
        INSERT INTO users (
            user_id, email, first_name, last_name,
            password_hash, provider_id, phone_number,
            provider, country
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
    `

	_, err := conn.Exec(ctx,
		query,
		usr.UserID,
		usr.Email,
		usr.FirstName,
		usr.LastName,
		usr.PasswordHash,
		usr.ProviderID,
		usr.PhoneNumber,
		usr.Provider,
		usr.Country,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_email_uq":
				return users.ErrEmailAlreadyExists
			case "users_user_id_uq":
				return users.ErrUserAlreadyExists
			case "users_provider_id_uq":
				return users.ErrProviderIDExists
			case "users_phone_number_uq":
				return users.ErrPhoneNumberExists
			case "provider_fields_chk":
				return users.ErrProviderFieldsCheck
			}
		}
		return fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	return nil
}

func (us *UserStore) FindUserByEmail(ctx context.Context, email string) (*users.User, error) {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
        SELECT user_id, email, first_name, last_name,
               password_hash, provider_id, phone_number,
               provider, country
        FROM users
        WHERE email = $1
    `

	var u users.User
	err := conn.QueryRow(ctx, query, email).Scan(
		&u.UserID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.PasswordHash,
		&u.ProviderID,
		&u.PhoneNumber,
		&u.Provider,
		&u.Country,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, users.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}
	return &u, nil
}

func (us *UserStore) FindUserPhoneNumber(ctx context.Context, phone string) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `SELECT 1 FROM users WHERE phone_number = $1`
	var ph int
	err := conn.QueryRow(ctx, query, phone).Scan(&ph)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("database error: %w", err)
	}
	return users.ErrPhoneNumberExists
}

func (us *UserStore) UpdateUser(ctx context.Context, usr *users.User) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
        UPDATE users
        SET email = $1,
            first_name = $2,
            last_name = $3,
            password_hash = $4,
            provider_id = $5,
            phone_number = $6,
            provider = $7,
            country = $8,
            updated_at = now()
        WHERE user_id = $9
    `
	_, err := conn.Exec(ctx,
		query,
		usr.Email,
		usr.FirstName,
		usr.LastName,
		usr.PasswordHash,
		usr.ProviderID,
		usr.PhoneNumber,
		usr.Provider,
		usr.Country,
		usr.UserID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_email_uq":
				return users.ErrEmailAlreadyExists
			case "users_user_id_uq":
				return users.ErrUserAlreadyExists
			case "users_provider_id_uq":
				return users.ErrProviderIDExists
			case "users_phone_number_uq":
				return users.ErrPhoneNumberExists
			case "provider_fields_chk":
				return users.ErrProviderFieldsCheck
			}
		}
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (us *UserStore) FindUserByID(ctx context.Context, userID uuid.UUID) (*users.User, error) {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
        SELECT user_id, email, first_name, last_name,
               password_hash, provider_id, phone_number,
               provider, country
        FROM users
        WHERE user_id = $1
    `

	var u users.User
	err := conn.QueryRow(ctx, query, userID).Scan(
		&u.UserID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.PasswordHash,
		&u.ProviderID,
		&u.PhoneNumber,
		&u.Provider,
		&u.Country,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, users.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}
	return &u, nil
}

func (us *UserStore) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
        UPDATE users
        SET is_verified = TRUE,
            updated_at = now()
        WHERE user_id = $1
    `
	cmdTag, err := conn.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed update users verification: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return users.ErrUserNotFound
	}

	return nil
}

func (us *UserStore) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
        UPDATE users
        SET password_hash = $1,
            updated_at = now()
        WHERE user_id = $2
    `
	cmdTag, err := conn.Exec(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return users.ErrUserNotFound
	}

	return nil
}
