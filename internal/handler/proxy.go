package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/FallPrediction/api-gateway/internal/metric"
	"github.com/FallPrediction/api-gateway/internal/response"
	"github.com/FallPrediction/api-gateway/internal/upstream"
)

var _ Handler = (*Gateway)(nil)

type Gateway struct{}

func (h *Gateway) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := upstream.GetUpstream(r.Context())
		if u == nil {
			http.NotFound(w, r)
		} else {
			url, _ := url.Parse(u.Upstream)
			proxy := httputil.NewSingleHostReverseProxy(url)
			rw := response.NewResponseWriter(w)
			start := time.Now()

			proxy.ServeHTTP(rw, r)

			metric.RequestDuration.WithLabelValues(u.Name, r.URL.Path).Observe(time.Since(start).Seconds())
			metric.RequestsTotal.WithLabelValues(u.Name, r.Method, strconv.Itoa(rw.StatusCode), r.URL.Path).Inc()
		}
	})
}

func NewGateway() Gateway {
	return Gateway{}
}
