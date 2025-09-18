package helper

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type cfg struct {
	Routes []Route `yaml:"origins"`
}

type Route struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Origin   string `yaml:"origin"`
	Auth     bool   `yaml:"auth"`
	RateLimt struct {
		Rate float64 `yaml:"rate"`
		Max  int     `yaml:"max"`
	} `yaml:"rate_limit"`
}

var Routes map[string]Route

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg cfg
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	Routes = make(map[string]Route)
	for _, route := range cfg.Routes {
		Routes[route.Path] = route
	}
}

func GetOrigin(url string) (Route, bool) {
	serviceName := ""
	parts := strings.SplitN(url, "/", 1)
	if len(parts) > 0 {
		serviceName = parts[0]
	}
	route, ok := Routes[serviceName]
	return route, ok
}
