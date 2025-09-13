package helper

import "golang.org/x/time/rate"

func NewRateLimiter() *rate.Limiter {
	return rate.NewLimiter(10, 10)
}
