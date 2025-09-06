package middleware

import (
	"api-gateway/logger"
	"net/http"

	"go.uber.org/zap"
)

var _ Middleware = new(Log)

type Log struct {
	baseMiddleware
}

func (m *Log) Handle() http.Handler {
	logger := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"",
			zap.String("type", "request"),
			zap.String("method", r.Method),
			zap.String("uri", r.URL.Path),
			zap.String("host", r.Host),
		)
		m.next.ServeHTTP(w, r)
	})
}

func NewLog() Log {
	return Log{}
}
