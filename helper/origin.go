package helper

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Route struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
	Auth     bool   `yaml:"auth"`
}

var Origins map[string]Route

func init() {
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

	Origins = make(map[string]Route)
	for _, route := range cfg.Routes {
		Origins[route.Path] = route
	}
}

func GetOrigin(url string) (Route, bool) {
	serviceName := ""
	parts := strings.SplitN(url, "/", 1)
	if len(parts) > 0 {
		serviceName = parts[0]
	}
	route, ok := Origins[serviceName]
	return route, ok
}
