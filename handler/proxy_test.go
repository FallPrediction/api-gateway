package handler_test

import (
	"api-gateway/handler"
	"api-gateway/helper"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGateway_Handle(t *testing.T) {
	t.Run("proxies the request to its configured upstream", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/orders", r.URL.Path)
			assert.Equal(t, "application/json", r.Header["Accept"][0])
			assert.Equal(t, "status[]=1&status[]=2&order_by=date", r.URL.RawQuery)
			// mock reqsponse
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "{\"error\":\"\",\"data\":[{\"id\":123,\"price\":234}]}")
		}))
		defer ts.Close()

		req := httptest.NewRequest("GET", "/orders?status[]=1&status[]=2&order_by=date", nil)
		req.Header.Set("Accept", "application/json")
		req = req.WithContext(helper.WithUpstream(req.Context(), &helper.Upstream{Name: "order", Upstream: ts.URL}))
		w := httptest.NewRecorder()
		handler := handler.NewGateway()
		handler.Handle().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		body, _ := io.ReadAll(w.Result().Body)
		w.Result().Body.Close()
		assert.Equal(t, []byte("{\"error\":\"\",\"data\":[{\"id\":123,\"price\":234}]}"), body)
	})

	t.Run("returns not found when no upstream is in the request context", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler := handler.NewGateway()

		handler.Handle().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/orders", nil))
		_, _ = io.Copy(io.Discard, w.Result().Body)
		w.Result().Body.Close()
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}
