package auth

import (
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

func (us *UserService) UserSignin(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)

	var usrReq UserSignInReq
	if err := base.ReadJSON(r, &usrReq); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.signin").Str("status", "failure").
			Str("req_id", reqID).Msg("failed to read JSON request")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	err := errs.NewValidate(usrReq)
	if err != nil {
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	signinUser, err := us.users.FindUserByEmail(r.Context(), usrReq.Email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			us.log.Warn().Str("event", "user.signin").Str("status", "failure").
				Str("req_id", reqID).Msg("user not found")
			base.WriteJSONError(w, errs.InvalidArgument, err)
			return
		}
		us.log.Error().
			Err(err).Str("event", "user.signin").Str("status", "failure").
			Str("req_id", reqID).Msg("database error fetching user")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	if err := authz.ComparePassword(signinUser.PasswordHash, []byte(usrReq.Password)); err != nil {
		if errors.Is(err, authz.ErrInvalidPassword) {
			us.log.Warn().Str("event", "user.signin").Str("status", "failure").
				Str("req_id", reqID).Msg("invalid credentials")
			base.WriteJSONError(w, errs.Unauthenticated, errors.New("invalid credentials"))
			return
		}
		us.log.Error().
			Err(err).Str("event", "user.signin").Str("status", "failure").
			Str("req_id", reqID).Msg("error comparing password")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	clientIP := base.GetClientIP(r)
	userAgent := r.UserAgent()

	accessTokenData := authz.NewJWTData(signinUser.UserID, signinUser.GetRole(), 30*time.Minute)
	accessToken, accesstokenPayload, err := us.auth.GenerateToken(accessTokenData)
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.signin").Str("status", "failure").
			Str("req_id", reqID).Msg("failed to generate access token")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	refreshTokenData := authz.NewJWTData(signinUser.UserID, signinUser.GetRole(), 24*time.Hour)
	refreshToken, refreshtokenPayload, err := us.auth.GenerateToken(refreshTokenData)
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.signin").Str("status", "failure").
			Str("req_id", reqID).Msg("failed to generate refresh token")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	hashRefresh := sha256.Sum256([]byte(refreshToken))
	err = us.users.CreateSession(r.Context(), &users.Session{
		ID:           uuid.MustParse(refreshtokenPayload.ID),
		UserID:       signinUser.UserID,
		RefreshToken: hashRefresh[:],
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		IsBlocked:    false,
		ExpiresAt:    refreshtokenPayload.ExpiresAt.Time,
	})
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.signin").Str("status", "failure").
			Str("req_id", reqID).Msg("failed to create user session")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	userData := UserSignInResp{
		TokenType:             "bearer",
		User:                  toUserResp(*signinUser),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accesstokenPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshtokenPayload.ExpiresAt.Time,
	}

	us.log.Info().
		Str("event", "user.signin").Str("status", "success").Str("req_id", reqID).
		Str("user_id", signinUser.UserID.String()).Msg("user signed in successfully")

	if err := base.WriteJSON(w, http.StatusOK, userData); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.signin").Str("status", "failure").
			Str("req_id", reqID).Msg("failed to write JSON response")
	}
}

func (us *UserService) SignOut(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)
	userPayload, ok := r.Context().Value(middleware.AuthContextPayloadKey).(*authz.Payload)
	if !ok {
		us.log.Warn().Str("event", "user.sign_out").
			Str("req_id", reqID).Msg("unauthenticated request to sign out")
		base.WriteJSONError(w, errs.Unauthenticated, errors.New("unauthorized"))
		return
	}

	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		us.log.Warn().
			Str("event", "user.sign_out").Str("req_id", reqID).
			Str("user_id", userPayload.UserID.String()).Msg("missing refresh token on sign out")
		base.WriteJSONError(w, errs.InvalidArgument, errors.New("missing refresh token"))
		return
	}

	shaToken := sha256.Sum256([]byte(refreshToken))

	err := us.users.BlockSession(r.Context(), shaToken[:])
	if err != nil {
		switch {
		case errors.Is(err, users.ErrSessionNotFound):
			us.log.Warn().
				Err(err).Str("event", "user.sign_out").Str("req_id", reqID).
				Str("user_id", userPayload.UserID.String()).Msg("session not found during sign out")
			base.WriteJSONError(w, errs.InvalidArgument, errors.New("invalid or expired session"))

		default:
			us.log.Error().
				Err(err).Str("event", "user.sign_out").Str("req_id", reqID).
				Str("user_id", userPayload.UserID.String()).Msg("unexpected error during sign out")
			base.WriteJSONInternalError(w, errs.Internal)
		}
		return
	}

	us.log.Info().
		Str("event", "user.sign_out").Str("status", "success").Str("req_id", reqID).
		Str("user_id", userPayload.UserID.String()).Msg("user signed out successfully")

	base.WriteJSON(w, http.StatusNoContent, nil)
}

func (us *UserService) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.RequestIdKey).(string)

	var req RenewAccessTokenReq
	if err := base.ReadJSON(r, &req); err != nil {
		us.log.Error().
			Err(err).Str("event", "user.renew_access_token").
			Str("req_id", reqID).Msg("invalid request body")
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	err := errs.NewValidate(req)
	if err != nil {
		base.WriteJSONError(w, errs.InvalidArgument, err)
		return
	}

	token, err := us.auth.VerifyToken(req.RefreshToken)
	if err != nil {
		base.WriteJSONError(w, errs.InvalidArgument, errors.New("invalid token or expired token"))
		return
	}

	session, err := us.users.GetSession(r.Context(), uuid.MustParse(token.ID))
	if err != nil {
		switch {
		case errors.Is(err, users.ErrTokenNotFound),
			errors.Is(err, users.ErrTokenExpired):
			us.log.Warn().
				Err(err).Str("event", "user.renew_access_token").Str("req_id", reqID).
				Msg("refresh token invalid or expired")
			base.WriteJSONError(w, errs.InvalidArgument, errors.New("refresh token invalid or expired"))
		case errors.Is(err, users.ErrUserNotFound):
			us.log.Warn().
				Err(err).Str("req_id", reqID).
				Str("event", "user.renew_access_token").Msg("user not found for refresh token")
			base.WriteJSONError(w, errs.InvalidArgument, err)
		default:
			us.log.Error().
				Err(err).Str("event", "user.renew_access_token").Str("req_id", reqID).
				Str("user_id", token.UserID.String()).Msg("unexpected database error during token renewal")
			base.WriteJSONInternalError(w, errs.Internal)
		}
		return
	}

	if session.IsBlocked {
		base.WriteJSONError(w, errs.InvalidArgument, errors.New("session is blocked"))
		return
	}

	if token.UserID != session.UserID {
		base.WriteJSONError(w, errs.InvalidArgument, errors.New("incorrect session user"))
		return
	}

	tokenData := authz.NewJWTData(session.UserID, token.RoleID, time.Minute*15)
	accessToken, accessPayload, err := us.auth.GenerateToken(tokenData)
	if err != nil {
		us.log.Error().
			Err(err).Str("event", "user.renew_access_token").Str("req_id", reqID).
			Str("user_id", session.UserID.String()).Msg("failed to create new access token")
		base.WriteJSONInternalError(w, errs.Internal)
		return
	}

	us.log.Info().
		Str("event", "user.renew_access_token").Str("status", "success").Str("req_id", reqID).
		Str("user_id", session.UserID.String()).Msg("access token renewed successfully")

	infoData := RenewAccessTokenResp{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}

	base.WriteJSON(w, http.StatusOK, infoData)
}
