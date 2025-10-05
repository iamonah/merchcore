package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"net/http"
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/users"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/base"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"
	"github.com/IamOnah/storefronthq/internal/sdk/middleware"

	"github.com/google/uuid"
)

func (us *UserService) RegisterUser(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)

	var usrReq UserCreateReq
	if err := base.ReadJSON(r, &usrReq); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.register").Str("status", "failure").
			Str("req_id", reqID).Msg("failed to read JSON request")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	err := errs.NewValidate(usrReq)
	if err != nil {
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	newUser, err := users.NewUser(toCreateUser(usrReq))
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.register").Str("status", "failure").
			Str("req_id", reqID).Msg("invalid user payload")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	if err := us.users.CreateUser(r.Context(), &newUser); err != nil {
		switch {
		case errors.Is(err, users.ErrDatabase):
			us.log.Error().
				Err(err).Str("event", "user.register").Str("status", "failure").
				Str("req_id", reqID).Msg("database error creating user")
			base.WriteJSONInternalError(w, errs.Internal)
		default:
			us.log.Warn().
				Err(err).Str("event", "user.register").Str("status", "failure").
				Str("req_id", reqID).Msg("failed to create user (invalid input or duplicate)")
			base.WriteJSONError(w, errs.InvalidArgument, err)
		}
		return
	}

	token, err := users.GenerateOTP(newUser.UserID, 90*time.Second, users.ActivationToken)
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.register").Str("status", "failure").
			Str("req_id", reqID).Str("user_id", newUser.UserID.String()).
			Msg("failed to generate OTP")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if err := us.users.CreateToken(r.Context(), token); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.register").Str("status", "failure").
			Str("req_id", reqID).Str("user_id", newUser.UserID.String()).Msg("database error creating OTP")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if err := us.job.WelcomeEmailJob(newUser.FirstName, token.Plaintext, newUser.UserID); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.register").Str("status", "failure").
			Str("user_id", newUser.UserID.String()).Msg("failed to enqueue welcome email job")
	} else {
		us.log.Info().
			Str("event", "user.register").Str("status", "success").
			Str("user_id", newUser.UserID.String()).Msg("welcome email job enqueued")
	}

	us.log.Info().
		Str("event", "user.register").Str("status", "success").Str("req_id", reqID).
		Str("user_id", newUser.UserID.String()).Msg("user successfully registered")

	infoData := UserCreateResp{
		Message: "user successfully registered, verify your account using the OTP sent to your email",
	}
	if err := base.WriteJSON(w, http.StatusCreated, infoData); err != nil {
		us.log.Error().Str("event", "user.register").
			Str("status", "failure").Str("user_id", newUser.UserID.String()).
			Err(err).Msg("failed to write JSON response")
	}
}

func (us *UserService) ActivateUser(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)

	var token TokenReq
	if err := base.ReadJSON(r, &token); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.activate").
			Str("req_id", reqID).Msg("failed to read JSON request")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	err := errs.NewValidate(token)
	if err != nil {
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}
	shaOTP := sha256.Sum256([]byte(token.Token))
	var userID uuid.UUID

	err = us.trx.WithTransaction(r.Context(), func(ctx context.Context) error {
		id, err := us.users.GetUserIDByToken(ctx, shaOTP[:], string(users.ActivationToken))
		if err != nil {
			return err
		}
		if err := us.users.DeleteToken(ctx, shaOTP[:], string(users.ActivationToken)); err != nil {
			return err
		}
		if err := us.users.VerifyUser(ctx, id); err != nil {
			return err
		}
		userID = id
		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, users.ErrDatabase):
			us.log.Error().
				Err(err).Str("event", "user.activate").
				Str("req_id", reqID).Msg("database error during user activation")
			base.WriteJSONInternalError(w, errs.Internal)

		case errors.Is(err, users.ErrTokenNotFound),
			errors.Is(err, users.ErrTokenExpired),
			errors.Is(err, users.ErrUserNotFound):
			us.log.Warn().
				Err(err).Str("event", "user.activate").Str("req_id", reqID).
				Str("user_id", userID.String()).Msg("invalid or expired OTP")
			base.WriteJSONError(w, errs.InvalidArgument, err)

		default:
			us.log.Error().
				Err(err).Str("event", "user.activate").Str("req_id", reqID).
				Msg("unexpected error activating user")
			base.WriteJSONInternalError(w, errs.Internal)
		}
		return
	}

	us.log.Info().
		Str("event", "user.activate").Str("status", "success").Str("req_id", reqID).
		Str("user_id", userID.String()).Msg("user successfully activated")

	base.WriteJSON(w, http.StatusNoContent, nil)
}

func (us *UserService) ResendVerificationToken(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)
	userPayload, ok := r.Context().Value(middleware.AuthContextPayloadKey).(*authz.Payload)
	if !ok {
		us.log.Warn().
			Str("event", "user.resend_verification").Str("req_id", reqID).
			Msg("unauthenticated request to get activation token")
		base.WriteJSONError(w, errs.Unauthenticated, errors.New("unauthorized"))
		return
	}

	newUser, err := us.users.FindUserByID(r.Context(), userPayload.UserID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			base.WriteJSONError(w, errs.InvalidArgument, err)
			return
		}
		us.log.Error().
			Err(err).Str("event", "user.resend_verification").Str("req_id", reqID).
			Str("user_id", userPayload.UserID.String()).Msg("failed to fetch user from DB")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if newUser.IsVerified {
		base.WriteJSONError(w, errs.InvalidArgument, errors.New("user is already verified"))
		return
	}

	token, err := users.GenerateOTP(userPayload.UserID, 90*time.Second, users.ActivationToken)
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.resend_verification").Str("req_id", reqID).
			Str("user_id", userPayload.UserID.String()).Msg("failed to generate OTP")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if err := us.users.CreateToken(r.Context(), token); err != nil {
		us.log.Error().
			Err(err).
			Str("event", "otp.save").Str("req_id", reqID).
			Str("user_id", userPayload.UserID.String()).Msg("database error creating OTP")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if err := us.job.VerificationEmailJob(newUser.FirstName, token.Plaintext, newUser.UserID); err != nil {
		us.log.Error().
			Err(err).
			Str("event", "email.job").Str("req_id", reqID).
			Str("user_id", newUser.UserID.String()).Msg("failed to enqueue verification email")
	} else {
		us.log.Info().
			Str("event", "email.job").Str("req_id", reqID).
			Str("user_id", newUser.UserID.String()).Msg("verification email job enqueued")
	}

	err = base.WriteJSON(w, http.StatusOK, map[string]string{"message": "success, OTP sent to your email"})
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.resend_verification").Str("req_id", reqID).
			Str("user_id", newUser.UserID.String()).Msg("failed to write JSON response")
	}
}
