package helper

import (
	"context"
	"os"

	"github.com/goccy/go-yaml"
)

type cfg struct {
	Upstreams []Upstream `yaml:"upstreams"`
}

type Upstream struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
	Auth     bool   `yaml:"auth"`
	RateLimt struct {
		Rate float64 `yaml:"rate"`
		Max  int     `yaml:"max"`
	} `yaml:"rate_limit"`
	Cors struct {
		AllowedMethods      []string `yaml:"allowed_methods"`
		AllowedOrigin       string   `yaml:"allowed_origin"`
		AllowedHeaders      []string `yaml:"allowed_headers"`
		ExposedHeaders      []string `yaml:"exposed_headers"`
		MaxAge              int      `yaml:"max_age"`
		SupportsCredentials bool     `yaml:"supports_credentials"`
	} `yaml:"cors"`
}

var Upstreams map[string]Upstream

type upstreamContextKey struct{}

var upstreamCtxKey upstreamContextKey

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

func WithUpstream(ctx context.Context, upstream *Upstream) context.Context {
	return context.WithValue(ctx, upstreamCtxKey, upstream)
}

func GetUpstream(requestContext context.Context) *Upstream {
	if v := requestContext.Value(upstreamCtxKey); v != nil {
		upstream, ok := v.(*Upstream)
		if ok {
			return upstream
		}
	}
	return nil
}
