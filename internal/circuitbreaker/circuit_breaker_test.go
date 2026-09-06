package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/FallPrediction/api-gateway/internal/circuitbreaker"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	circuitBreaker := circuitbreaker.NewCircuitBreaker(2, time.Second)
	assert.Equal(t, circuitBreaker.GetState(), circuitbreaker.Closed)
	circuitBreaker.RecordFailure()
	circuitBreaker.RecordFailure()
	assert.Equal(t, circuitBreaker.GetState(), circuitbreaker.Open)
	time.Sleep(time.Second)
	assert.Equal(t, circuitBreaker.GetState(), circuitbreaker.HalfOpen)
	assert.True(t, circuitBreaker.TrialCallStart())
	circuitBreaker.TrialCallOver()
	circuitBreaker.Reset()
	assert.Equal(t, circuitBreaker.GetState(), circuitbreaker.Closed)
}
