package helper

import (
	"golang.org/x/time/rate"
)

func NewRateLimiters(upstreams map[string]Upstream) map[string]*rate.Limiter {
	limiters := make(map[string]*rate.Limiter)
	for _, upstream := range upstreams {
		if upstream.RateLimt.Rate > 0 && upstream.RateLimt.Max > 0 {
			limiters[upstream.Name] = rate.NewLimiter(rate.Limit(upstream.RateLimt.Rate), upstream.RateLimt.Max)
		} else {
			limiters[upstream.Name] = NewDefaultRateLimiter()
		}
	}
	return limiters
}

func NewDefaultRateLimiter() *rate.Limiter {
	return rate.NewLimiter(500, 10000)
}
