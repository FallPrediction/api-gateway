package middleware

import "net/http"

type Middleware interface {
	Handle() http.Handler
	SetNext(next http.Handler)
}

type baseMiddleware struct {
	next http.Handler
}

func (m *baseMiddleware) SetNext(next http.Handler) {
	m.next = next
}
