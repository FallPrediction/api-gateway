package helper

import (
	"sync/atomic"
	"time"
)

type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	failureCount     atomic.Uint32
	failureThreshold uint32
	resetTimeout     time.Duration
	lastFailureTime  atomic.Pointer[time.Time]
	halfOpenInFlight atomic.Bool
}

func (c *CircuitBreaker) GetState() CircuitBreakerState {
	if c.failureCount.Load() < c.failureThreshold {
		return Closed
	}

	lastFailureTime := c.lastFailureTime.Load()
	if lastFailureTime != nil && time.Since(*lastFailureTime) > c.resetTimeout {
		return HalfOpen
	}

	return Open
}

func (c *CircuitBreaker) Reset() {
	c.failureCount.Store(0)
	c.lastFailureTime.Store(nil)
	c.halfOpenInFlight.Store(false)
}

func (c *CircuitBreaker) RecordFailure() {
	c.failureCount.Add(1)
	now := time.Now()
	c.lastFailureTime.Store(&now)
}

func (c *CircuitBreaker) TrialCallStart() bool {
	return c.halfOpenInFlight.CompareAndSwap(false, true)
}

func (c *CircuitBreaker) TrialCallOver() {
	c.halfOpenInFlight.Store(false)
}

func NewNewCircuitBreakers() map[string]*CircuitBreaker {
	limiters := make(map[string]*CircuitBreaker)
	for _, upstream := range Upstreams {
		limiters[upstream.Name] = NewCircuitBreaker()
	}
	return limiters
}

func NewCircuitBreaker() *CircuitBreaker {
	resetTimeout, _ := time.ParseDuration("30s")
	return &CircuitBreaker{
		failureThreshold: 5,
		resetTimeout:     resetTimeout,
	}
}
