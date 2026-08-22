package helper

import (
	"context"
	"os"

	"github.com/goccy/go-yaml"
)

type cfg struct {
	Upstreams []Upstream `yaml:"upstreams"`
}

type RateLimit struct {
	Rate float64 `yaml:"rate"`
	Max  int     `yaml:"max"`
}

type Cors struct {
	AllowedMethods      []string `yaml:"allowed_methods"`
	AllowedOrigin       string   `yaml:"allowed_origin"`
	AllowedHeaders      []string `yaml:"allowed_headers"`
	ExposedHeaders      []string `yaml:"exposed_headers"`
	MaxAge              int      `yaml:"max_age"`
	SupportsCredentials bool     `yaml:"supports_credentials"`
}

type Upstream struct {
	Name     string    `yaml:"name"`
	Path     string    `yaml:"path"`
	Upstream string    `yaml:"upstream"`
	Auth     bool      `yaml:"auth"`
	RateLimt RateLimit `yaml:"rate_limit"`
	Cors     Cors      `yaml:"cors"`
}

type upstreamContextKey struct{}

var upstreamCtxKey upstreamContextKey

func LoadConfig(configPath string) (map[string]Upstream, error) {
	upstreams := make(map[string]Upstream)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return upstreams, err
	}

	var cfg cfg
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return upstreams, err
	}

	for _, upstream := range cfg.Upstreams {
		upstreams[upstream.Path] = upstream
	}
	return upstreams, err
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
