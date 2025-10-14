package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/IamOnah/storefronthq/internal/sdk/base"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"

	"github.com/rs/zerolog"
)

func RecoverPanic(log *zerolog.Logger) base.Middleware {
	return func(next base.HTTPHandlerWithErr) base.HTTPHandlerWithErr {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error().
						Interface("panic", rec).
						Bytes("stack", debug.Stack()).
						Msg("panic recovered")
					err = errs.New(errs.Internal, errors.New("server temporarily unavailable"))
				}
			}()
			return next(w, r)
		}
	}
}
