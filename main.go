package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Route struct {
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
}

func getServices() map[string]Route {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg struct {
        Routes []Route
    }
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	services := make(map[string]Route)
	for _, route := range cfg.Routes {
		services[route.Path] = route
	}
	return services
}

func getServiceName(url string) string {
	parts := strings.SplitN(url, "/", 1)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func main() {
	services := getServices()
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		serviceName := getServiceName(req.URL.Path)
		route, ok := services[serviceName]
		if ok {
			url, _ := url.Parse(route.Upstream)
			proxy := httputil.NewSingleHostReverseProxy(url)
			proxy.ServeHTTP(rw, req)
		} else {
			http.NotFound(rw, req)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
