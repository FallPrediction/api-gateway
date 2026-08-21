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
	clean := path.Clean("/" + strings.TrimLeft(requestPath, "/"))
	if clean == "/" {
		return ""
	}
	segment := strings.SplitN(strings.TrimPrefix(clean, "/"), "/", 2)[0]
	if segment == "" {
		return ""
	}
	return "/" + segment
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
