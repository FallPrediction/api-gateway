package middleware

import (
	"net/http"

	"github.com/FallPrediction/api-gateway/internal/circuitbreaker"
	"github.com/FallPrediction/api-gateway/internal/metric"
	"github.com/FallPrediction/api-gateway/internal/response"
	"github.com/FallPrediction/api-gateway/internal/upstream"
)

var _ Middleware = (*CircuitBreaker)(nil)

type CircuitBreaker struct {
	baseMiddleware
	CircuitBreakers map[string]*circuitbreaker.CircuitBreaker
}

func (m *CircuitBreaker) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := upstream.GetUpstream(r.Context()).Name
		circuitBreaker := m.CircuitBreakers[u]
		if circuitBreaker == nil {
			m.next.ServeHTTP(w, r)
			return
		}

		state := circuitBreaker.GetState()
		if state == circuitbreaker.Open {
			metric.CircuitBreakerOpen.WithLabelValues(u).Inc()
			response.JSONResponse(
				w,
				http.StatusServiceUnavailable,
				map[string]string{},
				map[string]string{"msg": "Server is overloaded. Please try later."},
			)
			return
		}

		if state == circuitbreaker.HalfOpen && !circuitBreaker.TrialCallStart() {
			metric.CircuitBreakerOpen.WithLabelValues(u).Inc()
			response.JSONResponse(
				w,
				http.StatusServiceUnavailable,
				map[string]string{},
				map[string]string{"msg": "Server is overloaded. Please try later."},
			)
			return
		}

		if state == circuitbreaker.HalfOpen {
			defer circuitBreaker.TrialCallOver()
		}

		rw := response.NewResponseWriter(w)
		m.next.ServeHTTP(rw, r)

		if rw.StatusCode >= http.StatusInternalServerError {
			circuitBreaker.RecordFailure()
			return
		}

		circuitBreaker.Reset()
	})
}

func NewCircuitBreaker(circuitBreakers map[string]*circuitbreaker.CircuitBreaker) CircuitBreaker {
	return CircuitBreaker{
		CircuitBreakers: circuitBreakers,
	}
}
