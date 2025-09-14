package middleware

import (
	"api-gateway/helper"
	"net/http"

	"golang.org/x/time/rate"
)

type RateLimit struct {
	limiter *rate.Limiter
	baseMiddleware
}

var _ Middleware = new(RateLimit)

func (m *RateLimit) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			helper.GatewayRequestTotal.WithLabelValues("reject", r.URL.Path).Inc()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		helper.GatewayRequestTotal.WithLabelValues("accept", r.URL.Path).Inc()
		m.next.ServeHTTP(w, r)
	})
}

func NewRateLimit(limiter *rate.Limiter) RateLimit {
	return RateLimit{limiter: limiter}
}
