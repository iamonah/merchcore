package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/base"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"
)

type AuthKey string

const (
	AuthHeaderAuthorization AuthKey = "Authorization"
	AuthTypeBearer          AuthKey = "Bearer"
	AuthContextPayloadKey   AuthKey = "authorization_payload"
)

func AuthBearer(authMaker authz.TokenMaker) base.Middlware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(string(AuthHeaderAuthorization))
			if authHeader == "" {
				base.WriteJSONError(w, errs.Unauthenticated, errors.New("no auth header"))
				return
			}

			authParams := strings.Fields(authHeader)
			if len(authParams) != 2 || AuthKey(authParams[0]) != AuthTypeBearer {
				base.WriteJSONError(w, errs.Unauthenticated, errors.New("malformed auth header"))
				return
			}

			accessToken := authParams[1]
			payload, err := authMaker.VerifyToken(accessToken)
			if err != nil {
				base.WriteJSONError(w, errs.Unauthenticated, errors.New("invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), AuthContextPayloadKey, payload)
			next(w, r.WithContext(ctx))
		}
	}
}
