package middleware_test

import (
	"api-gateway/helper"
	"api-gateway/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpstream_Handle(t *testing.T) {
	upstream := helper.Upstream{Name: "order"}
	m := middleware.NewUpstream(map[string]helper.Upstream{"/order": upstream})

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
				gotUpstream := helper.GetUpstream(r.Context())
				if tt.returnNotFound == false && assert.NotNil(t, gotUpstream, "downstream request should contain upstream") {
					assert.Equal(t, upstream, *gotUpstream)
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
