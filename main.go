package main

import (
	"api-gateway/handler"
	"api-gateway/middleware"
	"net/http"
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
	logMiddleware := middleware.NewLog()
	recoveryMiddleware := middleware.NewRecover()
	authenticateMiddleware := middleware.NewAuthenticate()
	http.Handle("/", getHandler(&proxy, &logMiddleware, &recoveryMiddleware, &authenticateMiddleware))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
