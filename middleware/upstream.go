package middleware

import (
	"api-gateway/helper"
	"net/http"
	"path"
	"strings"
)

var _ Middleware = (*Upstream)(nil)

type Upstream struct {
	upstreams map[string]helper.Upstream
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
		upstream, ok := m.upstreams[upstreamKey]
		if ok {
			r = r.WithContext(helper.WithUpstream(r.Context(), &upstream))
			m.next.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

func NewUpstream(upstreams map[string]helper.Upstream) Upstream {
	return Upstream{upstreams: upstreams}
}
