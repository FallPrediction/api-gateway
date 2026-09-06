package middleware

import (
	"net/http"
	"path"
	"strings"

	"github.com/FallPrediction/api-gateway/internal/upstream"
)

var _ Middleware = (*Upstream)(nil)

type Upstream struct {
	upstreams map[string]upstream.Upstream
	baseMiddleware
}

func upstreamKeyFromRequestPath(requestPath string) string {
	clean := strings.TrimLeft(path.Clean(requestPath), "./")
	if clean == "" {
		return ""
	}
	return "/" + strings.SplitN(clean, "/", 2)[0]
}

func (m *Upstream) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamKey := upstreamKeyFromRequestPath(r.URL.Path)
		u, ok := m.upstreams[upstreamKey]
		if ok {
			r = r.WithContext(upstream.WithUpstream(r.Context(), &u))
			m.next.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

func NewUpstream(upstreams map[string]upstream.Upstream) Upstream {
	return Upstream{upstreams: upstreams}
}
