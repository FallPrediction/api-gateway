package middleware

import (
	"math"
	"net/http"
	"strconv"

	"github.com/FallPrediction/api-gateway/internal/metric"
	"github.com/FallPrediction/api-gateway/internal/ratelimit"
	"github.com/FallPrediction/api-gateway/internal/response"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"golang.org/x/time/rate"
)

type RateLimit struct {
	limiters       map[string]*rate.Limiter
	defaultLimiter *rate.Limiter
	baseMiddleware
}

var _ Middleware = (*RateLimit)(nil)

func (m *RateLimit) getLimiter(u *upstream.Upstream) *rate.Limiter {
	if u == nil {
		return m.defaultLimiter
	}
	if limiter := m.limiters[u.Name]; limiter != nil {
		return limiter
	}
	return m.defaultLimiter
}

func (m *RateLimit) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := upstream.GetUpstream(r.Context())
		reservation := m.getLimiter(u).Reserve()
		if !reservation.OK() {
			reservation.Cancel()
			metric.GatewayRequestTotal.WithLabelValues(u.Name, "reject", r.URL.Path).Inc()
			response.JSONResponse(
				w,
				http.StatusTooManyRequests,
				map[string]string{},
				map[string]string{
					"msg": "Request tokens exceed the Limiter's burst size",
				},
			)
			return
		} else if reservation.Delay() > 0 {
			reservation.Cancel()
			metric.GatewayRequestTotal.WithLabelValues(u.Name, "reject", r.URL.Path).Inc()
			response.JSONResponse(
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
		metric.GatewayRequestTotal.WithLabelValues(u.Name, "accept", r.URL.Path).Inc()
		m.next.ServeHTTP(w, r)
	})
}

func NewRateLimit(limiters map[string]*rate.Limiter) RateLimit {
	return RateLimit{
		limiters:       limiters,
		defaultLimiter: ratelimit.NewDefaultRateLimiter(),
	}
}
