package auth

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/iamonah/merchcore/internal/sdk/base"
	"github.com/iamonah/merchcore/internal/sdk/errs"
)

func (us *UserService) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	var req ForgotPasswordReq
	if err := base.ReadJSON(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	email, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errs.New(errs.InvalidArgument, errors.New("invalid email"))
	}
	user, token, err := us.users.ForgetPassword(r.Context(), email)
	if err != nil {
		return errs.Newf(errs.Internal, "forgetpassword: %v", err)
	}
	if token == nil && user == nil {
		info := map[string]string{"message": "reset link sent if account exists"}
		if err := base.WriteJSON(w, http.StatusOK, info); err != nil {
			return errs.Newf(errs.Internal, "writejson: %s", err)
		}
		return nil
	}

	if err := us.job.PasswordResetEmailJob(user.Email.Address, token.Plaintext, user.UserID); err != nil {
		return errs.Newf(errs.Internal, "passwordresetjob: userID[%+v]: %v", user.UserID.ID(), err)
	}

	data := map[string]string{"message": "reset link sent if account exists"}
	if err := base.WriteJSON(w, http.StatusOK, data); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}
	return nil
}

// followup for forgotpassword
func (us *UserService) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	reqID, err := base.GetReqIDCTX(r)
	if err != nil {
		return errs.Newf(errs.Internal, "getreqidCTX: %s", err)
	}

	var req ResetPasswordReq
	if err := base.ReadJSON(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if req.NewPassword != req.ConfirmPassword {
		return errs.New(errs.InvalidArgument, errors.New("passwords do not match"))
	}

	userID, err := us.users.PasswordReset(r.Context(), req.NewPassword, req.Token)
	if err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, derr)
		}
		return errs.Newf(errs.Internal, "passwordreset: %s", err)
	}

	us.log.Info().
		Str("event", "user.reset_password").Str("status", "success").
		Str("req_id", reqID).Str("user_id", userID.String()).
		Msg("password successfully reset")

	data := map[string]string{"message": "password successfully reset"}
	if err := base.WriteJSON(w, http.StatusOK, data); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}

	return nil
}

func (us *UserService) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	reqID, err := base.GetReqIDCTX(r)
	if err != nil {
		return errs.Newf(errs.Internal, "getreqidCTX: %s", err)
	}
	pl, err := base.GetJWTPayloadCTX(r)
	if err != nil {
		return errs.New(errs.Unauthenticated, errors.New("unauthorized"))
	}

	var req ChangePasswordReq
	if err := base.ReadJSON(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if req.NewPassword != req.ConfirmPassword {
		return errs.New(errs.InvalidArgument, errors.New("passwords do not match"))
	}

	user, err := us.users.ChangePassword(r.Context(), pl.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		return errs.Newf(errs.Internal, "changepassword: userID[%v]: %s", pl.UserID, err)
	}

	us.log.Info().
		Str("event", "user.change_password").Str("status", "success").Str("req_id", reqID).
		Str("user_id", user.UserID.String()).Msg("password successfully changed")

	data := map[string]string{"message": "password successfully changed"}
	if err := base.WriteJSON(w, http.StatusOK, data); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}

	return nil
}
