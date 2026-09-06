package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FallPrediction/api-gateway/internal/handler"
	"github.com/FallPrediction/api-gateway/internal/middleware"
	"github.com/FallPrediction/api-gateway/internal/testutil"
	"github.com/FallPrediction/api-gateway/internal/upstream"

	"github.com/stretchr/testify/assert"
)

func Test_getHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/order", r.URL.Path)
		assert.Equal(t, "application/json", r.Header["Accept"][0])
		assert.Equal(t, "status[]=1&status[]=2&order_by=date", r.URL.RawQuery)
		// mock reqsponse
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "{\"error\":\"\",\"data\":[{\"id\":123,\"price\":234}]}")
	}))
	defer ts.Close()

	proxy := handler.NewGateway()
	u := middleware.NewUpstream(map[string]upstream.Upstream{"/order": {Name: "order", Path: "/order", Upstream: ts.URL, Auth: true}})
	authenticate := middleware.NewAuthenticate()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", ts.URL+"/order?status[]=1&status[]=2&order_by=date", nil)
	token, _ := testutil.CreateToken("order", time.Now().Add(time.Minute))
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	handler := getHandler(&proxy, &u, &authenticate)
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	body, _ := io.ReadAll(w.Result().Body)
	w.Result().Body.Close()
	assert.Equal(t, []byte("{\"error\":\"\",\"data\":[{\"id\":123,\"price\":234}]}"), body)
}
