package helper

import (
	"golang.org/x/time/rate"
)

func NewRateLimiters() map[string]*rate.Limiter {
	limiters := make(map[string]*rate.Limiter)
	for _, origin := range Origins {
		if origin.RateLimt.Rate > 0 && origin.RateLimt.Max > 0 {
			limiters[origin.Name] = rate.NewLimiter(rate.Limit(origin.RateLimt.Rate), origin.RateLimt.Max)
		} else {
			limiters[origin.Name] = NewDefaultRateLimiter()
		}
	}
	return limiters
}

func NewDefaultRateLimiter() *rate.Limiter {
	return rate.NewLimiter(500, 10000)
}
