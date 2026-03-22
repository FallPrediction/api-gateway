package main

import (
	"api-gateway/handler"
	"api-gateway/helper"
	"api-gateway/middleware"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getHandler(handler handler.Handler, middlewares ...middleware.Middleware) http.Handler {
	if len(middlewares) == 0 {
		return handler.Handle()
	}
	for i := 0; i < len(middlewares)-1; i++ {
		middlewares[i].SetNext(middlewares[i+1].Handle())
	}
	middlewares[len(middlewares)-1].SetNext(handler.Handle())
	return middlewares[0].Handle()
}

func main() {
	proxy := handler.NewGateway()
	recoveryMiddleware := middleware.NewRecover()
	upstreamMiddleware := middleware.NewUpstream()
	rateLimitMiddleware := middleware.NewRateLimit(helper.NewRateLimiters())
	logMiddleware := middleware.NewLog()
	authenticateMiddleware := middleware.NewAuthenticate()
	http.Handle("/", getHandler(
		&proxy,
		&recoveryMiddleware,
		&upstreamMiddleware,
		&rateLimitMiddleware,
		&logMiddleware,
		&authenticateMiddleware,
	))
	http.Handle("/metrics", promhttp.HandlerFor(helper.NewRegistry(), promhttp.HandlerOpts{}))
	http.Handle("/health", handler.NewHealth().Handle())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
