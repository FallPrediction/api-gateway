package ratelimit

import (
	"github.com/FallPrediction/api-gateway/internal/upstream"
	"golang.org/x/time/rate"
)

func NewRateLimiters(upstreams map[string]upstream.Upstream) map[string]*rate.Limiter {
	limiters := make(map[string]*rate.Limiter)
	for _, u := range upstreams {
		if u.RateLimt.Rate > 0 && u.RateLimt.Max > 0 {
			limiters[u.Name] = rate.NewLimiter(rate.Limit(u.RateLimt.Rate), u.RateLimt.Max)
		} else {
			limiters[u.Name] = NewDefaultRateLimiter()
		}
	}
	return limiters
}

func NewDefaultRateLimiter() *rate.Limiter {
	return rate.NewLimiter(500, 10000)
}
