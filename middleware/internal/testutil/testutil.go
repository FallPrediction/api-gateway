package testutil

import (
	"api-gateway/helper"
	"net/http"
	"net/http/httptest"
)

func CreateRequestWithUpstream(upstream *helper.Upstream) *http.Request {
	req := httptest.NewRequest("GET", "/"+upstream.Name, nil)
	req = req.WithContext(helper.WithUpstream(req.Context(), upstream))
	return req
}

func MockResponse(statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	})
}
