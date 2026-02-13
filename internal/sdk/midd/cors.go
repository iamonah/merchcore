package midd

import (
	"net/http"
	"strings"

	"github.com/iamonah/merchcore/internal/config"
)

const origin = "Origin"

func EnableCors(cfg config.Config) Middleware {
	return func(next HTTPHandlerWithErr) HTTPHandlerWithErr {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Methods")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			origin := r.Header.Get(origin)

			if origin != "" {
				for _, o := range cfg.Server.CORSAllowedOrigins {
					if strings.TrimSpace(o) == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)

						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, PATCH, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
							w.Header().Set("Access-Control-Max-Age", "300")
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
