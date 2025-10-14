package userdb

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/IamOnah/storefronthq/internal/domain/shared/contact"
	"github.com/IamOnah/storefronthq/internal/domain/shared/role"
	"github.com/IamOnah/storefronthq/internal/domain/users"
	"github.com/IamOnah/storefronthq/internal/infra/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userdb struct {
	conn database.DBTX
}

func Newuserdb(conn *pgxpool.Pool) *userdb {
	return &userdb{
		conn: conn,
	}
}
func (us *userdb) CreateUser(ctx context.Context, usr *users.User) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
		INSERT INTO users (user_id, email, first_name, last_name,
			password_hash, provider_id, phone_number,
			provider, country, number_of_store, is_store_created, role) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING created_at, updated_at
	`

	err := conn.QueryRow(ctx,
		query,
		usr.UserID,
		usr.Email.Address,
		usr.FirstName,
		usr.LastName,
		usr.PasswordHash,
		usr.ProviderID,
		usr.Contact.Number,
		usr.Provider,
		usr.Contact.Country,
		usr.NumOfStore,
		usr.IsStoreCreated,
		usr.Role.String(),
	).Scan(&usr.CreatedAt, &usr.UpdatedAT)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_email_uq":
				return users.ErrEmailAlreadyExists
			case "users_user_id_uq":
				return users.ErrUserIDConflict
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

func (us *userdb) GetUserByEmail(ctx context.Context, emailValue string) (*users.User, error) {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
		SELECT user_id, email, first_name, last_name, password_hash,
			provider_id, phone_number, provider, country,
			created_at, updated_at, is_verified, deleted_at,
			role, is_store_created, number_of_store AS num_of_store
		FROM users
		WHERE email = $1;        
    `

	var (
		emailStr string
		phoneNum string
		country  string
		roleStr  string
		u        users.User
	)
	err := conn.QueryRow(ctx, query, emailValue).Scan(
		&u.UserID,
		&emailStr,
		&u.FirstName,
		&u.LastName,
		&u.PasswordHash,
		&u.ProviderID,
		&phoneNum,
		&u.Provider,
		&country,
		&u.CreatedAt,
		&u.UpdatedAT,
		&u.IsVerified,
		&u.DeletedAt,
		&roleStr,
		&u.IsStoreCreated,
		&u.NumOfStore,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, users.ErrUserNotFound
		}
		return nil, fmt.Errorf("getuserbyemail: %w: %w", users.ErrDatabase, err)
	}
	email, err := mail.ParseAddress(emailStr)
	if err != nil {
		return nil, fmt.Errorf("getuserbyemail: %w: invalid stored email %s: %w", users.ErrDatabase, emailStr, err)
	}

	u.Email = email

	contact := contact.NewContact(phoneNum, country)
	err = contact.ValidateContact()
	if err != nil {
		return nil, fmt.Errorf("getuserbyemail: %w: invalid stored contact %s,%s: %w", users.ErrDatabase, phoneNum, country, err)
	}

	u.Role, err = role.Parse(roleStr)
	if err != nil {
		return nil, fmt.Errorf("getuserbyemail: invalid stored role %q: %w", roleStr, err)
	}

	u.Contact = contact
	return &u, nil
}

func (us *userdb) GetUserPhoneNumber(ctx context.Context, phone string) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `SELECT 1 FROM users WHERE phone_number = $1`
	var ph int
	err := conn.QueryRow(ctx, query, phone).Scan(&ph)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("geruserbyphonenumber: %w: %w", users.ErrDatabase, err)
	}
	return users.ErrPhoneNumberExists
}

func (us *userdb) UpdateUser(ctx context.Context, usr *users.User) error {
	conn := database.GetTXFromContext(ctx, us.conn)
	query := `
		UPDATE users
		SET
			email = $1,
			first_name = $2,
			last_name = $3,
			password_hash = $4,
			provider_id = $5,
			phone_number = $6,
			provider = $7,
			country = $8,
			role = $9,
			number_of_store = $10,
			is_store_created = $11,
			updated_at = now()
		WHERE user_id = $12
	`
	_, err := conn.Exec(ctx,
		query,
		usr.Email.Address,
		usr.FirstName,
		usr.LastName,
		usr.PasswordHash,
		usr.ProviderID,
		usr.Contact.Number,
		usr.Provider,
		usr.Contact.Country,
		usr.Role.String(),
		usr.NumOfStore,
		usr.IsStoreCreated,
		usr.UserID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_email_uq":
				return users.ErrEmailAlreadyExists
			case "users_user_id_uq":
				return users.ErrUserIDConflict
			case "users_provider_id_uq":
				return users.ErrProviderIDExists
			case "users_phone_number_uq":
				return users.ErrPhoneNumberExists
			case "provider_fields_chk":
				return users.ErrProviderFieldsCheck
			default:
				return fmt.Errorf("unhandled db constraint: %s (%s)", pgErr.ConstraintName, pgErr.Message)
			}
		}
		return fmt.Errorf("updateyuser: %w: %w", users.ErrDatabase, err)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return users.ErrUserNotFound
	}

	return nil
}

func (us *userdb) GetUserByID(ctx context.Context, userID uuid.UUID) (*users.User, error) {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
		SELECT user_id, email, first_name, last_name, password_hash,
			   provider_id, phone_number, provider, country,
			   created_at, updated_at, is_verified, deleted_at,
			   role, is_store_created, number_of_store AS num_of_store
		FROM users
		WHERE user_id = $1;
	`

	var (
		emailStr string
		phoneNum string
		country  string
		roleStr  string
		u        users.User
	)

	err := conn.QueryRow(ctx, query, userID).Scan(
		&u.UserID,
		&emailStr,
		&u.FirstName,
		&u.LastName,
		&u.PasswordHash,
		&u.ProviderID,
		&phoneNum,
		&u.Provider,
		&country,
		&u.CreatedAt,
		&u.UpdatedAT,
		&u.IsVerified,
		&u.DeletedAt,
		&roleStr,
		&u.IsStoreCreated,
		&u.NumOfStore,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, users.ErrUserNotFound
		}
		return nil, fmt.Errorf("getuserbyid: %w: %w", users.ErrDatabase, err)
	}

	email, err := mail.ParseAddress(emailStr)
	if err != nil {
		return nil, fmt.Errorf("getuserbyid: %w: invalid stored email %s: %w", users.ErrDatabase, emailStr, err)
	}
	u.Email = email

	contact := contact.NewContact(phoneNum, country)
	if err := contact.ValidateContact(); err != nil {
		return nil, fmt.Errorf("getuserbyid: %w: invalid stored contact %s,%s: %w", users.ErrDatabase, phoneNum, country, err)
	}
	u.Contact = contact

	u.Role, err = role.Parse(roleStr)
	if err != nil {
		return nil, fmt.Errorf("getuserbyid: invalid stored role %q: %w", roleStr, err)
	}

	return &u, nil
}


func (us *userdb) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	query := `
        UPDATE users
        SET is_verified = TRUE,
            updated_at = now()
        WHERE user_id = $1
    `
	cmdTag, err := conn.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return users.ErrUserNotFound
	}

	return nil
}

func (us *userdb) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
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
