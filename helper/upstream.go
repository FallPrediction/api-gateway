package helper

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type cfg struct {
	Upstreams []Upstream `yaml:"upstreams"`
}

type Upstream struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Upstream   string `yaml:"upstream"`
	Auth     bool   `yaml:"auth"`
	RateLimt struct {
		Rate float64 `yaml:"rate"`
		Max  int     `yaml:"max"`
	} `yaml:"rate_limit"`
}

var Upstreams map[string]Upstream

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg cfg
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	Upstreams = make(map[string]Upstream)
	for _, upstream := range cfg.Upstreams {
		Upstreams[upstream.Path] = upstream
	}
}

func GetUpstream(url string) (Upstream, bool) {
	serviceName := ""
	parts := strings.SplitN(url, "/", 1)
	if len(parts) > 0 {
		serviceName = parts[0]
	}
	upstream, ok := Upstreams[serviceName]
	return upstream, ok
}

func GetUpstreamName(url string) string {
	service := "Unknown"
	upstream, ok := GetUpstream(url)
	if ok {
		service = upstream.Name
	}
	return service
}
