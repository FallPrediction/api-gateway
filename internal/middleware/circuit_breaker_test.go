package middleware_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FallPrediction/api-gateway/internal/circuitbreaker"
	"github.com/FallPrediction/api-gateway/internal/middleware"
	"github.com/FallPrediction/api-gateway/internal/testutil"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_Handle(t *testing.T) {
	failureThreshold := 2
	m := middleware.NewCircuitBreaker(map[string]*circuitbreaker.CircuitBreaker{
		"order": circuitbreaker.NewCircuitBreaker(failureThreshold, time.Second),
	})
	m.SetNext(testutil.MockResponse(http.StatusInternalServerError))

	u := &upstream.Upstream{Name: "order"}
	for range failureThreshold {
		w := httptest.NewRecorder()
		m.Handle().ServeHTTP(w, testutil.CreateRequestWithUpstream(u))
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode, "The last few calls returned 500 status code.")
	}

	w := httptest.NewRecorder()
	m.Handle().ServeHTTP(w, testutil.CreateRequestWithUpstream(u))
	assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode, "After failureThreshold failed calls, the next call returns 503 status code.")
	respBody := struct {
		Msg string
	}{}
	response, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(response, &respBody)
	assert.Equal(t, "Server is overloaded. Please try later.", respBody.Msg, "After failureThreshold failed calls, the next call returns the error message.")
	w.Result().Body.Close()
	time.Sleep(time.Second)

	w = httptest.NewRecorder()
	m.SetNext(testutil.MockResponse(http.StatusOK))
	_, _ = io.Copy(io.Discard, w.Result().Body)
	w.Result().Body.Close()
	m.Handle().ServeHTTP(w, testutil.CreateRequestWithUpstream(u))
	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Call again one second later to trigger a trial call and the status code of response is 200.")
}
