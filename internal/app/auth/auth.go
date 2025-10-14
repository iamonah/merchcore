package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/base"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"
	"github.com/IamOnah/storefronthq/internal/sdk/middleware"
)

func (us *UserService) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)

	var usrReq UserCreateReq
	if err := base.ReadJSON(r, &usrReq); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := errs.NewValidate(usrReq); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	newUser, token, err := us.users.CreateUser(r.Context(), toCreateUser(usrReq))
	if err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, err)
		}
		return errs.Newf(errs.Internal, "createuser: %s", err)
	}

	us.log.Info().
		Str("event", "user.register").
		Str("req_id", reqID).
		Str("user_id", newUser.UserID.String()).
		Msg("user registered")

	fmt.Println("big token", token.Plaintext)
	if err := us.job.WelcomeEmailJob(newUser.FirstName, token.Plaintext, newUser.UserID, newUser.Email.Address); err != nil {
		us.log.Error().
			Err(err).
			Str("event", "user.register").
			Str("user_id", newUser.UserID.String()).
			Msg("welcome email enqueue")
	} else {
		us.log.Info().
			Str("event", "user.register").
			Str("user_id", newUser.UserID.String()).
			Msg("welcome email enqueued")
	}
	info := UserCreateResp{Message: "succes, verify your email"}

	if err := base.WriteJSON(w, http.StatusCreated, info); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}
	return nil
}

func (us *UserService) ActivateUser(w http.ResponseWriter, r *http.Request) error {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)
	pl, ok := r.Context().Value(middleware.AuthContextPayloadKey).(*authz.Payload)
	if !ok {
		return errs.New(errs.Unauthenticated, errors.New("unauthorized"))
	}

	var t TokenReq
	if err := base.ReadJSON(r, &t); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	if err := errs.NewValidate(t); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := us.users.ActivateUser(r.Context(), pl.UserID, t.Token); err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, err)
		}
		return errs.Newf(errs.Internal, "activateuser: %s", err)
	}

	us.log.Info().
		Str("event", "user.activate").
		Str("req_id", reqID).
		Str("user_id", pl.UserID.String()).
		Msg("user activated")

	if err := base.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}
	return nil
}

func (us *UserService) ResendVerificationToken(w http.ResponseWriter, r *http.Request) error {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)
	userPayload, ok := r.Context().Value(middleware.AuthContextPayloadKey).(*authz.Payload)
	if !ok {
		return errs.New(errs.Unauthenticated, errors.New("unauthorized"))
	}

	user, token, err := us.users.ResendActivationToken(r.Context(), userPayload.UserID)
	if err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, err)
		}
		return errs.Newf(errs.Internal, "resendactivationtoken: userID[%s]: %s", userPayload.UserID, err)
	}

	if err := us.job.ResendVerificationTokenJob(user.FirstName, token.Plaintext, user.UserID); err != nil {
		return errs.Newf(errs.Internal, "resendactivation: enqueue job userID[%s]: %s", user.UserID, err)
	}

	us.log.Info().
		Str("event", "email.job").
		Str("req_id", reqID).
		Str("user_id", user.UserID.String()).
		Msg("verification email job enqueued")

	data := map[string]string{"message": "success, OTP sent to your email"}
	if err := base.WriteJSON(w, http.StatusOK, data); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}
	return nil
}
