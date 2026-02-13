package auth

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/iamonah/merchcore/internal/sdk/base"
	"github.com/iamonah/merchcore/internal/sdk/errs"
)

func (us *UserService) Authenticate(w http.ResponseWriter, r *http.Request) error {
	reqID, err := base.GetReqIDCTX(r)
	if err != nil {
		return errs.Newf(errs.Internal, "getreqidCTX: %s", err)
	}

	var req UserSignInReq
	if err := base.ReadJSON(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := errs.NewValidate(req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	email, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errs.New(errs.Unauthenticated, errors.New("invalid credentials"))
	}

	user, err := us.users.Authenticate(r.Context(), email, req.Password)
	if err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, derr)
		}
		return errs.Newf(errs.Internal, "authenticate: user[%s]: %s", email.Address, err)
	}

	clientIP := base.GetClientIP(r)
	userAgent := r.UserAgent()

	session, err := us.users.CreateSession(r.Context(), user, userAgent, clientIP)
	if err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, derr)
		}
		return errs.Newf(errs.Internal, "authenticate: user[%s]: %s", user.UserID, err)
	}

	resp := UserSignInResp{
		TokenType:             "bearer",
		User:                  toUserResp(user),
		AccessToken:           session.AccessToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
	}

	if err := base.WriteJSON(w, http.StatusOK, resp); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}

	us.log.Info().
		Str("event", "user.signin").
		Str("req_id", reqID).
		Str("user_id", user.UserID.String()).
		Msg("authenticate: success")

	return nil
}

func (us *UserService) SignOut(w http.ResponseWriter, r *http.Request) error {
	reqID, err := base.GetReqIDCTX(r)
	if err != nil {
		return errs.Newf(errs.Internal, "getreqidCTX: %s", err)
	}
	pl, err := base.GetJWTPayloadCTX(r)
	if err != nil {
		return errs.New(errs.Unauthenticated, errors.New("unauthorized"))
	}

	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		return errs.New(errs.InvalidArgument, errors.New("missing refresh token"))
	}

	if err := us.users.BlockSession(r.Context(), pl.UserID, refreshToken); err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, derr)
		}
		return errs.Newf(errs.Internal, "signout: user[%s]: %s", pl.UserID, err)
	}

	us.log.Info().
		Str("event", "user.signout").
		Str("req_id", reqID).
		Str("user_id", pl.UserID.String()).
		Msg("signout: success")

	if err := base.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}
	return nil
}

func (us *UserService) RenewAccessToken(w http.ResponseWriter, r *http.Request) error {
	reqID, err := base.GetReqIDCTX(r)
	if err != nil {
		return errs.Newf(errs.Internal, "getreqidCTX: %s", err)
	}
	pl, err := base.GetJWTPayloadCTX(r)
	if err != nil {
		return errs.New(errs.Unauthenticated, errors.New("unauthorized"))
	}
	token, err := us.users.RenewAccessToken(r.Context(), pl)
	if err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, derr)
		}
		return errs.Newf(errs.Internal, "renewaccess: user[%s]: %s", token.UserId, err)
	}

	resp := RenewAccessTokenResp{
		AccessToken:          token.AccessToken,
		AccessTokenExpiresAt: token.AcessExpiresAt,
	}

	us.log.Info().
		Str("event", "user.renew_access_token").
		Str("req_id", reqID).
		Str("user_id", token.UserId.String()).
		Msg("renew_token: success")

	if err := base.WriteJSON(w, http.StatusOK, resp); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}
	return nil
}
