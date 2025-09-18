package middleware

import (
	"api-gateway/helper"
	"net/http"

	"golang.org/x/time/rate"
)

type RateLimit struct {
	limiters       map[string]*rate.Limiter
	defaultLimiter *rate.Limiter
	baseMiddleware
}

var _ Middleware = new(RateLimit)

func (m *RateLimit) getLimiter(url string) *rate.Limiter {
	route, ok := helper.GetOrigin(url)
	if ok {
		return m.limiters[route.Name]
	}
	return m.defaultLimiter
}

func (m *RateLimit) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.getLimiter(r.URL.Path).Allow() {
			helper.GatewayRequestTotal.WithLabelValues("reject", r.URL.Path).Inc()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		helper.GatewayRequestTotal.WithLabelValues("accept", r.URL.Path).Inc()
		m.next.ServeHTTP(w, r)
	})
}

func NewRateLimit(limiters map[string]*rate.Limiter) RateLimit {
	return RateLimit{
		limiters:       limiters,
		defaultLimiter: helper.NewDefaultRateLimiter(),
	}
}
