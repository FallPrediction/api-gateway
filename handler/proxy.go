package handler

import (
	"api-gateway/helper"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var _ Handler = new(Gateway)

type Gateway struct{}

func (h *Gateway) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, ok := helper.GetServiceName(r.URL.Path)
		if ok {
			url, _ := url.Parse(route.Upstream)
			proxy := httputil.NewSingleHostReverseProxy(url)
			proxy.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

func NewGateway() Gateway {
	return Gateway{}
}
