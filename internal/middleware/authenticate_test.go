package middleware_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FallPrediction/api-gateway/internal/middleware"
	"github.com/FallPrediction/api-gateway/internal/testutil"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type CustomerInfo struct {
	Name string
	Kind string
}

type CustomClaimsExample struct {
	jwt.RegisteredClaims
	TokenType string
	CustomerInfo
}

func TestAuthenticate_Handle(t *testing.T) {
	t.Setenv("APP_KEY", "x02tiYd31fV9EUdNqhJ3a7jTSwXbhHCjgny6xli9pwQ=")
	m := middleware.NewAuthenticate()
	m.SetNext(testutil.MockResponse(http.StatusOK))

	tests := []struct {
		name       string
		req        *http.Request
		statusCode int
		repMsg     string
	}{
		{
			"Request without token.",
			func() *http.Request {
				return testutil.CreateRequestWithUpstream(&upstream.Upstream{
					Name: "order",
					Auth: true,
				})
			}(),
			http.StatusUnauthorized,
			"Unauthorized",
		},
		{
			"Invalid token.",
			func() *http.Request {
				req := testutil.CreateRequestWithUpstream(&upstream.Upstream{
					Name: "order",
					Auth: true,
				})
				req.Header.Set("authorization", "Bearer abc")
				return req
			}(),
			http.StatusUnauthorized,
			"Unauthorized",
		},
		{
			"Token expired",
			func() *http.Request {
				token, _ := testutil.CreateToken("order", time.Now().Add(-time.Minute))
				req := testutil.CreateRequestWithUpstream(&upstream.Upstream{
					Name: "order",
					Auth: true,
				})
				req.Header.Set("authorization", "Bearer "+token)
				return req
			}(),
			http.StatusForbidden,
			"Token expired",
		},
		{
			"Token without scope",
			func() *http.Request {
				token, _ := testutil.CreateToken("", time.Now().Add(time.Minute))
				req := testutil.CreateRequestWithUpstream(&upstream.Upstream{
					Name: "order",
					Auth: true,
				})
				req.Header.Set("authorization", "Bearer "+token)
				return req
			}(),
			http.StatusForbidden,
			"Forbidden",
		},
		{
			"Authenticate pass",
			func() *http.Request {
				token, _ := testutil.CreateToken("order", time.Now().Add(time.Minute))
				req := testutil.CreateRequestWithUpstream(&upstream.Upstream{
					Name: "order",
					Auth: true,
				})
				req.Header.Set("authorization", "Bearer "+token)
				return req
			}(),
			http.StatusOK,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			m.Handle().ServeHTTP(w, tt.req)
			assert.Equal(t, tt.statusCode, w.Result().StatusCode, tt.name)
			respBody := struct {
				Msg string
			}{}
			body, _ := io.ReadAll(w.Result().Body)
			w.Result().Body.Close()
			json.Unmarshal(body, &respBody)
			assert.Equal(t, tt.repMsg, respBody.Msg, tt.name)
		})
	}
}
