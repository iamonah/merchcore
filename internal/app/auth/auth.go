package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/users"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/base"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"
	"github.com/IamOnah/storefronthq/internal/sdk/middleware"

	"github.com/google/uuid"
)

func (us *UserService) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)

	var req ForgotPasswordReq
	if err := base.ReadJSON(r, &req); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.forgot_password").
			Str("req_id", reqID).Msg("invalid request body")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	user, err := us.users.FindUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			data := map[string]string{"message": "A reset link will be been sent if account exist"}
			base.WriteJSON(w, http.StatusOK, data)
			return
		}
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	token, err := users.GenerateToken(user.UserID, 15*time.Minute, users.PasswordReset)
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.forgot_password").
			Str("req_id", reqID).Msg("failed to generate reset token")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if err := us.users.CreateToken(r.Context(), token); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.forgot_password").Str("req_id", reqID).
			Msg("database failed to store reset token")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if err := us.job.PasswordResetEmailJob(user.Email.Address, token.Plaintext, user.UserID); err != nil {
		us.log.Error().
			Str("event", "user.forgot_password").Str("req_id", reqID).
			Err(err).Msg("failed to enqueue reset email")
	}
	outputData := map[string]string{"message": "If the account exists, a reset link has been sent"}
	base.WriteJSON(w, http.StatusOK, outputData)
}

// User coming from reset link
func (us *UserService) ResetPassword(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)

	var req ResetPasswordReq
	if err := base.ReadJSON(r, &req); err != nil {
		us.log.Error().Str("event", "user.reset_password").Str("req_id", reqID).
			Err(err).Msg("invalid request body")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		us.log.Warn().
			Str("event", "user.reset_password").Str("req_id", reqID).
			Msg("new password and confirm password do not match")
		base.WriteJSONError(w, errs.InvalidArgument, errors.New("passwords do not match"))
		return
	}

	newPassword, err := authz.HashPassword([]byte(req.NewPassword))
	if err != nil {
		us.log.Error().Str("event", "user.reset_password").Str("req_id", reqID).
			Err(err).Msg("failed to hash password")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	shaToken := sha256.Sum256([]byte(req.Token))
	var userID uuid.UUID

	err = us.trx.WithTransaction(r.Context(), func(ctx context.Context) error {
		id, err := us.users.GetUserIDByToken(ctx, shaToken[:], string(users.PasswordReset))
		if err != nil {
			return err
		}
		userID = id

		if err := us.users.DeleteToken(ctx, shaToken[:], string(users.PasswordReset)); err != nil {
			return err
		}

		if err := us.users.UpdatePassword(ctx, userID, newPassword); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, users.ErrTokenNotFound),
			errors.Is(err, users.ErrTokenExpired),
			errors.Is(err, users.ErrUserNotFound):
			us.log.Warn().
				Err(err).Str("event", "user.reset_password").Str("req_id", reqID).
				Str("user_id", userID.String()).Msg("invalid or expired reset token")
			base.WriteJSONError(w, errs.InvalidArgument, errors.New("reset link expired or invalid"))

		default:
			us.log.Error().
				Err(err).Str("event", "user.reset_password").Str("req_id", reqID).
				Str("user_id", userID.String()).Msg("unexpected error during password reset")
			base.WriteJSONInternalError(w, errs.Internal)
		}
		return
	}

	us.log.Info().
		Str("event", "user.reset_password").Str("status", "success").
		Str("req_id", reqID).Str("user_id", userID.String()).
		Msg("password successfully reset")

	base.WriteJSON(w, http.StatusOK, map[string]string{"message": "password successfully reset"})
}

func (us *UserService) ChangePassword(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)
	userPayload, ok := r.Context().Value(middleware.AuthContextPayloadKey).(*authz.Payload)
	if !ok {
		us.log.Warn().Str("event", "user.change_password").
			Str("req_id", reqID).Msg("unauthenticated request to change password")
		base.WriteJSONError(w, errs.Unauthenticated, errors.New("unauthorized"))
		return
	}

	var req ChangePasswordReq
	if err := base.ReadJSON(r, &req); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.change_password").
			Str("req_id", reqID).Msg("invalid request body")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		us.log.Warn().
			Str("event", "user.change_password").
			Str("req_id", reqID).Msg("new password and confirm password do not match")
		base.WriteJSONError(w, errs.InvalidArgument, errors.New("passwords do not match"))
		return
	}

	var userData *users.User
	err := us.trx.WithTransaction(r.Context(), func(ctx context.Context) error {
		// Get current hashed password
		userValue, err := us.users.FindUserByID(ctx, userPayload.UserID)
		if err != nil {
			return err
		}

		if err := authz.ComparePassword(userValue.PasswordHash, []byte(req.OldPassword)); err != nil {
			return authz.ErrInvalidPassword
		}

		newHash, err := authz.HashPassword([]byte(req.NewPassword))
		if err != nil {
			return fmt.Errorf("hash new password: %w", err)
		}

		userValue.PasswordHash = newHash
		if err := us.users.UpdateUser(ctx, userValue); err != nil {
			return err
		}

		userData = userValue
		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			us.log.Warn().
				Err(err).Str("event", "user.change_password").Str("req_id", reqID).
				Str("user_id", userData.UserID.String()).Msg("user not found ")
			base.WriteJSONError(w, errs.InvalidArgument, err)

		case errors.Is(err, authz.ErrInvalidPassword):
			us.log.Warn().
				Err(err).Str("event", "user.change_password").Str("req_id", reqID).
				Str("user_id", userData.UserID.String()).Msg("old password incorrect")
			base.WriteJSONError(w, errs.InvalidArgument, errors.New("old password incorrect"))

		default:
			us.log.Error().
				Err(err).Str("event", "user.change_password").Str("req_id", reqID).
				Str("user_id", userData.UserID.String()).Msg("unexpected error during password change")
			base.WriteJSONInternalError(w, errs.Internal)
		}
		return
	}

	us.log.Info().
		Str("event", "user.change_password").Str("status", "success").Str("req_id", reqID).
		Str("user_id", userData.UserID.String()).Msg("password successfully changed")

	err = base.WriteJSON(w, http.StatusOK, map[string]string{"message": "password successfully changed"})
	if err != nil {
		base.WriteJSONInternalError(w, errs.Internal)
	}
}
