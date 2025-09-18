package helper

import (
	"golang.org/x/time/rate"
)

func NewRateLimiters() map[string]*rate.Limiter {
	limiters := make(map[string]*rate.Limiter)
	for _, route := range Origins {
		if route.RateLimt.rate > 0 && route.RateLimt.max > 0 {
			limiters[route.Name] = rate.NewLimiter(rate.Limit(route.RateLimt.rate), route.RateLimt.max)
		} else {
			limiters[route.Name] = NewDefaultRateLimiter()
		}
	}
	return limiters
}

func NewDefaultRateLimiter() *rate.Limiter {
	return rate.NewLimiter(500, 10000)
}
