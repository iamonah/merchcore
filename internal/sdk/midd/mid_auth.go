package midd

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/iamonah/merchcore/internal/sdk/authz"
	"github.com/iamonah/merchcore/internal/sdk/errs"
)

type AuthKey string

const (
	AuthHeaderAuthorization AuthKey = "Authorization"
	AuthTypeBearer          AuthKey = "Bearer"
	AuthContextPayloadKey   AuthKey = "authorization_payload"
)

func AuthBearer(authMaker authz.TokenMaker) Middleware {
	return func(next HTTPHandlerWithErr) HTTPHandlerWithErr {
		return func(w http.ResponseWriter, r *http.Request) error {
			authHeader := r.Header.Get(string(AuthHeaderAuthorization))
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", "Bearer")
				return errs.New(errs.Unauthenticated, errors.New("missing authorization header"))
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || AuthKey(parts[0]) != AuthTypeBearer {
				w.Header().Set("WWW-Authenticate", "Bearer")
				return errs.New(errs.Unauthenticated, errors.New("malformed authorization header"))
			}

			payload, err := authMaker.VerifyToken(parts[1])
			if err != nil {
				w.Header().Set("WWW-Authenticate", "Bearer")
				return errs.New(errs.Unauthenticated, errors.New("invalid or expired token"))
			}

			ctx := context.WithValue(r.Context(), AuthContextPayloadKey, payload)
			return next(w, r.WithContext(ctx))
		}
	}
}

//use redis for session midd()
