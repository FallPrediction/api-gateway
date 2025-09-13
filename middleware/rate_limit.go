package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
)

type RateLimit struct {
	limiter *rate.Limiter
	baseMiddleware
}

var _ Middleware = new(RateLimit)

func (m *RateLimit) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		m.next.ServeHTTP(w, r)
	})
}

func NewRateLimit(limiter *rate.Limiter) RateLimit {
	return RateLimit{limiter: limiter}
}
