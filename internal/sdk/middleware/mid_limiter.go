package middleware

// import (
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/IamOnah/storefronthq/internal/sdk/base"
// 	"golang.org/x/time/rate"
// )

// type limiter struct {
// 	rps     float64
// 	burst   int
// 	enabled bool
// }

// func RateLimit(next http.Handler) http.Handler {
// 	//On initialization of code
// 	type client struct {
// 		limiter  *rate.Limiter
// 		lastSeen time.Time
// 	}

// 	var (
// 		mu      sync.Mutex
// 		clients = make(map[string]*client)
// 	)

// 	go func() {
// 		for {
// 			time.Sleep(time.Minute)
// 			mu.Lock()
// 			for ip, client := range clients {
// 				if time.Since(client.lastSeen) > 3*time.Minute {
// 					delete(clients, ip)
// 				}
// 			}
// 			mu.Unlock()
// 		}

// 	}()

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if app.config.limiter.enabled {
// 			ip := base.GetClientIP(r)
// 			mu.Lock()

// 			if _, found := clients[ip]; !found {
// 				clients[ip] = &client{
// 					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
// 				}
// 			}
// 			clients[ip].lastSeen = time.Now()

// 			if !clients[ip].limiter.Allow() {
// 				mu.Unlock()
// 				app.rateLimitExceededResponse(w, r)
// 				return
// 			}
// 			mu.Unlock()
// 		}
// 		next.ServeHTTP(w, r)
// 	})

// }
