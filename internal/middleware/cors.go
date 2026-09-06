package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/FallPrediction/api-gateway/internal/upstream"
)

var _ Middleware = (*Cors)(nil)

type Cors struct {
	baseMiddleware
}

func (m *Cors) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := upstream.GetUpstream(r.Context()).Cors

		if config.AllowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", config.AllowedOrigin)
			if config.AllowedOrigin != "*" && strings.Contains(config.AllowedOrigin, "*") {
				w.Header().Add("Vary", "Origin")
			}
		}

		if len(config.AllowedMethods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		}
		if len(config.AllowedHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		}
		if len(config.ExposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
		}
		if config.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
		}
		if config.SupportsCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" && r.Header.Get("Origin") != "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		m.next.ServeHTTP(w, r)
	})
}

func NewCors() Cors {
	return Cors{}
}
