package middleware

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/iamonah/merchcore/internal/config"
	"github.com/iamonah/merchcore/internal/sdk/base"
	"github.com/iamonah/merchcore/internal/sdk/errs"
	"golang.org/x/time/rate"
)

type limiter struct {
	rps     float64
	burst   int
	enabled bool
}

func RateLimit(cfg config.Config) base.Middleware {
	//On initialization of code
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}

	}()

	return func(next base.HTTPHandlerWithErr) base.HTTPHandlerWithErr {
		return func(w http.ResponseWriter, r *http.Request) error {
			ip := base.GetClientIP(r)
			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					//cfg.limiter.rps
					//cfg.limiter.burst
					limiter: rate.NewLimiter(rate.Limit(5), 8),
				}
			}
			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				return errs.New(errs.TooManyRequests, errors.New("rate limit exceeded"))
			}
			return next(w, r)
		}
	}
}
