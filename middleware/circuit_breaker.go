package middleware

import (
	"api-gateway/helper"
	"net/http"
)

var _ Middleware = (*CircuitBreaker)(nil)

type CircuitBreaker struct {
	baseMiddleware
	CircuitBreakers map[string]*helper.CircuitBreaker
}

func (m *CircuitBreaker) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstream := helper.GetUpstream(r.Context()).Name
		circuitBreaker := m.CircuitBreakers[upstream]
		if circuitBreaker == nil {
			m.next.ServeHTTP(w, r)
			return
		}

		state := circuitBreaker.GetState()
		if state == helper.Open {
			helper.CircuitBreakerOpen.WithLabelValues(upstream).Inc()
			helper.JSONResponse(
				w,
				http.StatusServiceUnavailable,
				map[string]string{},
				map[string]string{"msg": "Server is overloaded. Please try later."},
			)
			return
		}

		if state == helper.HalfOpen && !circuitBreaker.TrialCallStart() {
			helper.CircuitBreakerOpen.WithLabelValues(upstream).Inc()
			helper.JSONResponse(
				w,
				http.StatusServiceUnavailable,
				map[string]string{},
				map[string]string{"msg": "Server is overloaded. Please try later."},
			)
			return
		}

		if state == helper.HalfOpen {
			defer circuitBreaker.TrialCallOver()
		}

		rw := helper.NewResponseWriter(w)
		m.next.ServeHTTP(rw, r)

		if rw.StatusCode >= http.StatusInternalServerError {
			circuitBreaker.RecordFailure()
			return
		}

		circuitBreaker.Reset()
	})
}

func NewCircuitBreaker(circuitBreakers map[string]*helper.CircuitBreaker) CircuitBreaker {
	return CircuitBreaker{
		CircuitBreakers: circuitBreakers,
	}
}
