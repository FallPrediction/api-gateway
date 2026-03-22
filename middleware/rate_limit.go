package middleware

import (
	"api-gateway/helper"
	"math"
	"net/http"
	"strconv"

	"golang.org/x/time/rate"
)

type RateLimit struct {
	limiters       map[string]*rate.Limiter
	defaultLimiter *rate.Limiter
	baseMiddleware
}

var _ Middleware = (*RateLimit)(nil)

func (m *RateLimit) getLimiter(upstream *helper.Upstream) *rate.Limiter {
	if upstream == nil {
		return m.defaultLimiter
	}
	if limiter := m.limiters[upstream.Name]; limiter != nil {
		return limiter
	}
	return m.defaultLimiter
}

func (m *RateLimit) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstream := helper.GetUpstream(r.Context())
		reservation := m.getLimiter(upstream).Reserve()
		if !reservation.OK() {
			helper.GatewayRequestTotal.WithLabelValues(upstream.Name, "reject", r.URL.Path).Inc()
			helper.JSONResponse(
				w,
				http.StatusTooManyRequests,
				map[string]string{
					"Retry-After": strconv.Itoa(int(math.Ceil(reservation.Delay().Seconds()))),
				},
				map[string]string{
					"msg": "Too Many Requests",
				},
			)
			return
		}
		helper.GatewayRequestTotal.WithLabelValues(upstream.Name, "accept", r.URL.Path).Inc()
		m.next.ServeHTTP(w, r)
	})
}

func NewRateLimit(limiters map[string]*rate.Limiter) RateLimit {
	return RateLimit{
		limiters:       limiters,
		defaultLimiter: helper.NewDefaultRateLimiter(),
	}
}
