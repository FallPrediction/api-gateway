package handler

import (
	"api-gateway/helper"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

var _ Handler = new(Gateway)

type Gateway struct{}

func (h *Gateway) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, ok := helper.GetServiceName(r.URL.Path)
		if ok {
			url, _ := url.Parse(route.Upstream)
			proxy := httputil.NewSingleHostReverseProxy(url)
			rw := helper.NewResponseWriter(w)
			start := time.Now()

			proxy.ServeHTTP(w, r)

			helper.RequestDuration.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())
			helper.RequestsTotal.WithLabelValues(r.Method, strconv.Itoa(rw.StatusCode), r.URL.Path).Inc()
		} else {
			http.NotFound(w, r)
		}
	})
}

func NewGateway() Gateway {
	return Gateway{}
}
