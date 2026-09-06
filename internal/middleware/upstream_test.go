package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FallPrediction/api-gateway/internal/middleware"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"github.com/stretchr/testify/assert"
)

func TestUpstream_Handle(t *testing.T) {
	u := upstream.Upstream{Name: "order"}
	m := middleware.NewUpstream(map[string]upstream.Upstream{"/order": u})

	tests := []struct {
		name           string
		uri            string
		returnNotFound bool
	}{
		{
			"Root path no upstream",
			"/",
			true,
		},
		{
			"Unknown upstream",
			"/user",
			true,
		},
		{
			"Double slash with upstream",
			"//order",
			false,
		},
		{
			"Double slash with upstream",
			"/order/1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.SetNext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotUpstream := upstream.GetUpstream(r.Context())
				if tt.returnNotFound == false && assert.NotNil(t, gotUpstream, "downstream request should contain upstream") {
					assert.Equal(t, u, *gotUpstream)
				}
				w.WriteHeader(http.StatusOK)
			}))

			w := httptest.NewRecorder()
			m.Handle().ServeHTTP(w, httptest.NewRequest(http.MethodGet, tt.uri, nil))
			_, _ = io.Copy(io.Discard, w.Result().Body)
			w.Result().Body.Close()
			if tt.returnNotFound {
				assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
			} else {
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			}
		})
	}
}
