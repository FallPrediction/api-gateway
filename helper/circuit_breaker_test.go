package helper_test

import (
	"api-gateway/helper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	circuitBreaker := helper.NewCircuitBreaker(2, time.Second)
	assert.Equal(t, circuitBreaker.GetState(), helper.Closed)
	circuitBreaker.RecordFailure()
	circuitBreaker.RecordFailure()
	assert.Equal(t, circuitBreaker.GetState(), helper.Open)
	time.Sleep(time.Second)
	assert.Equal(t, circuitBreaker.GetState(), helper.HalfOpen)
	assert.True(t, circuitBreaker.TrialCallStart())
	circuitBreaker.TrialCallOver()
	circuitBreaker.Reset()
	assert.Equal(t, circuitBreaker.GetState(), helper.Closed)
}
