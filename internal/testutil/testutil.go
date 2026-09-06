package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/FallPrediction/api-gateway/internal/upstream"
)

func CreateRequestWithUpstream(u *upstream.Upstream) *http.Request {
	req := httptest.NewRequest("GET", "/"+u.Name, nil)
	req = req.WithContext(upstream.WithUpstream(req.Context(), u))
	return req
}

func MockResponse(statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	})
}
