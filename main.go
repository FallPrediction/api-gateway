package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
)

func main() {
    usersURL, _ := url.Parse("http://users.com")
    ordersURL, _ := url.Parse("http://orders.com")

    usersProxy := httputil.NewSingleHostReverseProxy(usersURL)
    ordersProxy := httputil.NewSingleHostReverseProxy(ordersURL)

    http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
        switch {
        case strings.HasPrefix(req.URL.Path, "/users"):
            usersProxy.ServeHTTP(rw, req)
        case strings.HasPrefix(req.URL.Path, "/orders"):
            ordersProxy.ServeHTTP(rw, req)
        default:
            http.NotFound(rw, req)
        }
    })

    log.Println("Proxy server is running on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
