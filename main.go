package main

import (
	"api-gateway/handler"
	"api-gateway/middleware"
	"net/http"
)

func main() {
	proxy := handler.NewGateway()
	logMiddleware := middleware.NewLog()
	logMiddleware.SetNext(proxy.Handle())
	http.Handle("/", logMiddleware.Handle())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
