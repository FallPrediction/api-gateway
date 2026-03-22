package handler

import (
	"api-gateway/helper"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

var _ Handler = (*Gateway)(nil)

type Gateway struct{}

func (h *Gateway) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstream := helper.GetUpstream(r.Context())
		if upstream == nil {
			http.NotFound(w, r)
		} else {
			url, _ := url.Parse(upstream.Upstream)
			proxy := httputil.NewSingleHostReverseProxy(url)
			rw := helper.NewResponseWriter(w)
			start := time.Now()

			proxy.ServeHTTP(w, r)

			helper.RequestDuration.WithLabelValues(upstream.Name, r.URL.Path).Observe(time.Since(start).Seconds())
			helper.RequestsTotal.WithLabelValues(upstream.Name, r.Method, strconv.Itoa(rw.StatusCode), r.URL.Path).Inc()
		}
	})
}

func NewGateway() Gateway {
	return Gateway{}
}
