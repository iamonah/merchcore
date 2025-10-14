package middleware

import (
	"net/http"
	"strings"

	"github.com/IamOnah/storefronthq/internal/config"
	"github.com/IamOnah/storefronthq/internal/sdk/base"
)

func EnableCors(cfg config.Config) base.Middleware {
	return func(next base.HTTPHandlerWithErr) base.HTTPHandlerWithErr {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			origin := r.Header.Get("Origin")

			if origin != "" {
				for _, o := range cfg.Server.CORSAllowedOrigins {
					if strings.TrimSpace(o) == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)

						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, PATCH, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
							w.Header().Set("Access-Control-Max-Age", "1800")
							w.WriteHeader(http.StatusNoContent)
							return nil
						}
						break
					}
				}
			}
			return next(w, r)
		}
	}
}
