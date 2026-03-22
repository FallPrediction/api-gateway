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

func (m *RateLimit) getLimiter(url string) *rate.Limiter {
	upstream, ok := helper.GetUpstream(url)
	if ok {
		return m.limiters[upstream.Name]
	}
	return m.defaultLimiter
}

func (m *RateLimit) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstream := helper.GetUpstreamName(r.URL.Path)
		reservation := m.getLimiter(r.URL.Path).Reserve()
		if !reservation.OK() {
			helper.GatewayRequestTotal.WithLabelValues(upstream, "reject", r.URL.Path).Inc()
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
		helper.GatewayRequestTotal.WithLabelValues(upstream, "accept", r.URL.Path).Inc()
		m.next.ServeHTTP(w, r)
	})
}

func NewRateLimit(limiters map[string]*rate.Limiter) RateLimit {
	return RateLimit{
		limiters:       limiters,
		defaultLimiter: helper.NewDefaultRateLimiter(),
	}
}
