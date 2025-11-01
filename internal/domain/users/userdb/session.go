package userdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	users "github.com/iamonah/merchcore/internal/domain/users"
	"github.com/iamonah/merchcore/internal/infra/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (us *userdb) CreateSession(ctx context.Context, s *users.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6,$7)
	`
	conn := database.GetTXFromContext(ctx, us.conn)

	_, err := conn.Exec(
		ctx,
		query,
		s.ID,
		s.UserID,
		s.RefreshToken,
		s.UserAgent,
		s.ClientIP,
		s.IsBlocked,
		s.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	return nil
}

func (us *userdb) CreateToken(ctx context.Context, otp *users.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`
	conn := database.GetTXFromContext(ctx, us.conn)

	_, err := conn.Exec(ctx, query, otp.TokenHash[:], otp.UserID, otp.Expiry, otp.Scope)
	if err != nil {
		return fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}
	return nil
}

func (us *userdb) GetUserIDByToken(ctx context.Context, hash []byte, scope string) (uuid.UUID, error) {
	const query = `
		SELECT t.user_id, t.expiry
		FROM tokens t
		WHERE t.hash = $1 AND t.scope = $2
	`

	conn := database.GetTXFromContext(ctx, us.conn)

	var userID uuid.UUID
	var expiry time.Time

	err := conn.QueryRow(ctx, query, hash, scope).Scan(&userID, &expiry)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, users.ErrUserNotFound
		}
		return uuid.Nil, fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	if time.Now().After(expiry) {
		return uuid.Nil, users.ErrTokenExpired
	}

	return userID, nil
}

func (us *userdb) DeleteToken(ctx context.Context, hash []byte, scope string) error {
	const query = `
		DELETE FROM tokens	WHERE hash = $1 AND scope = $2
	`
	conn := database.GetTXFromContext(ctx, us.conn)

	cmdTag, err := conn.Exec(ctx, query, hash, scope)
	if err != nil {
		return fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return users.ErrTokenNotFound
	}
	return nil
}

func (us *userdb) BlockSession(ctx context.Context, token []byte) error {
	conn := database.GetTXFromContext(ctx, us.conn)

	const q = `
		UPDATE sessions
		SET is_blocked = true
		WHERE refresh_token = $1 AND NOT is_blocked
	`

	res, err := conn.Exec(ctx, q, token)
	if err != nil {
		return fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	if res.RowsAffected() == 0 {
		return users.ErrSessionNotFound
	}

	return nil
}
func (us *userdb) GetSession(ctx context.Context, sessionID string) (*users.Session, error) {
	conn := database.GetTXFromContext(ctx, us.conn)

	const query = `
        SELECT id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
        FROM sessions
        WHERE id = $1
    `

	var s users.Session
	err := conn.QueryRow(ctx, query, sessionID).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshToken,
		&s.UserAgent,
		&s.ClientIP,
		&s.IsBlocked,
		&s.ExpiresAt,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, users.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: %w", users.ErrDatabase, err)
	}

	return &s, nil
}
