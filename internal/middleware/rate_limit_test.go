package middleware_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/FallPrediction/api-gateway/internal/middleware"
	"github.com/FallPrediction/api-gateway/internal/ratelimit"
	"github.com/FallPrediction/api-gateway/internal/testutil"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit_Handle(t *testing.T) {
	u := upstream.Upstream{
		Name: "order",
		RateLimt: upstream.RateLimit{
			Rate: 1,
			Max:  1,
		},
	}
	m := middleware.NewRateLimit(ratelimit.NewRateLimiters(map[string]upstream.Upstream{
		"/order": u,
	}))
	m.SetNext(testutil.MockResponse(http.StatusOK))
	req := testutil.CreateRequestWithUpstream(&u)

	w := httptest.NewRecorder()
	m.Handle().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "First request return 200")
	_, _ = io.Copy(io.Discard, w.Result().Body)
	w.Result().Body.Close()

	w = httptest.NewRecorder()
	m.Handle().ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Result().StatusCode, "Second request return 429")
	respBody := struct {
		Msg string
	}{}
	body, _ := io.ReadAll(w.Result().Body)
	w.Result().Body.Close()
	json.Unmarshal(body, &respBody)
	assert.Equal(t, "Too Many Requests", respBody.Msg, "Second request message should be Too Many Requests")
	retryAfter, _ := strconv.Atoi(w.Result().Header["Retry-After"][0])
	assert.LessOrEqual(t, retryAfter, 1, "Second request should wait less or equal 1 second to send request again")
}
