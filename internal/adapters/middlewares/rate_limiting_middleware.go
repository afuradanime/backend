package middlewares

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// Rate limiting by IP
type IPRateLimiter struct {
	limiters sync.Map // concurrent safe map
	Rps      float64  // How many token refreshes an ip gets per second (e.g: Burst = 10 means new token every 100ms)
	Burst    int      // How many requests can be made "instantly"
}

// rate.Limiter uses https://en.wikipedia.org/wiki/Token_bucket
func (l *IPRateLimiter) get(ip string) *rate.Limiter {
	v, _ := l.limiters.LoadOrStore(ip, rate.NewLimiter(rate.Limit(l.Rps), l.Burst))
	return v.(*rate.Limiter)
}

func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !l.get(ip).Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
