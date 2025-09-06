package main

import (
	"api-gateway/handler"
	"net/http"
)

func main() {
	proxy := handler.NewGateway()
	http.Handle("/", proxy.Handle())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
