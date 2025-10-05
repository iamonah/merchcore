package middleware

import (
	"net/http"
	"strings"

	"github.com/IamOnah/storefronthq/internal/config"
	"github.com/IamOnah/storefronthq/internal/sdk/base"
)

func EnableCors(cfg config.Config) base.Middlware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			origin := r.Header.Get("Origin")

			if origin != "" {
				for i := range cfg.Server.CORSAllowedOrigins {
					cleanOrigins := strings.TrimSpace(cfg.Server.CORSAllowedOrigins[i])
					if origin == cleanOrigins {
						w.Header().Set("Access-Control-Allow-Origin", origin)

						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST,GET, PUT, PATCH, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
							w.Header().Set("Access-Control-Max-Age", "1800")
							w.WriteHeader(http.StatusNoContent)
							return
						}
						break
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
