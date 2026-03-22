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
		upstream, ok := helper.GetUpstream(r.URL.Path)
		if ok {
			url, _ := url.Parse(upstream.Upstream)
			proxy := httputil.NewSingleHostReverseProxy(url)
			rw := helper.NewResponseWriter(w)
			start := time.Now()

			proxy.ServeHTTP(w, r)

			helper.RequestDuration.WithLabelValues(helper.GetUpstreamName(r.URL.Path), r.URL.Path).Observe(time.Since(start).Seconds())
			helper.RequestsTotal.WithLabelValues(helper.GetUpstreamName(r.URL.Path), r.Method, strconv.Itoa(rw.StatusCode), r.URL.Path).Inc()
		} else {
			http.NotFound(w, r)
		}
	})
}

func NewGateway() Gateway {
	return Gateway{}
}
