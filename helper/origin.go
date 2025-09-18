package helper

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type cfg struct {
	Origins []Origin `yaml:"origins"`
}

type Origin struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Origin   string `yaml:"origin"`
	Auth     bool   `yaml:"auth"`
	RateLimt struct {
		Rate float64 `yaml:"rate"`
		Max  int     `yaml:"max"`
	} `yaml:"rate_limit"`
}

var Origins map[string]Origin

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg cfg
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	Origins = make(map[string]Origin)
	for _, origin := range cfg.Origins {
		Origins[origin.Path] = origin
	}
}

func GetOrigin(url string) (Origin, bool) {
	serviceName := ""
	parts := strings.SplitN(url, "/", 1)
	if len(parts) > 0 {
		serviceName = parts[0]
	}
	origin, ok := Origins[serviceName]
	return origin, ok
}
