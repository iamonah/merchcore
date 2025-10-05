package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/IamOnah/storefronthq/internal/sdk/base"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"

	"github.com/rs/zerolog"
)

func RecoverPanic(log *zerolog.Logger) base.Middlware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					log.Error().
						Any("panic", err).
						Bytes("stack", debug.Stack()).
						Msg("panic recovered")

					w.Header().Set("Connection", "close") //successfull close the connection
					base.WriteJSONError(w, errs.Internal, errors.New("server temporarily unavailable"))
				}
			}()
			next.ServeHTTP(w, r)
		}
	}
}
