package middleware

import (
	"api-gateway/helper"
	"api-gateway/logger"
	"io"
	"math/rand/v2"
	"net/http"

	"go.uber.org/zap"
)

var _ Middleware = new(Log)

type Log struct {
	baseMiddleware
}

func (m *Log) getTraceId() string {
	resultBytes := [5]byte{}
	character := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range len(resultBytes) {
		resultBytes[i] = character[rand.IntN(len(character))]
	}
	return string(resultBytes[:])
}

func (m *Log) Handle() http.Handler {
	logger := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := helper.NewResponseWriter(w)
		body, _ := io.ReadAll(r.Body)
		traceId := m.getTraceId()
		logger.Info(
			"",
			zap.String("trace_id", traceId),
			zap.String("type", "request"),
			zap.String("method", r.Method),
			zap.String("uri", r.URL.Path),
			zap.String("host", r.Host),
			zap.Object("header", (*helper.ZapHeader)(&r.Header)),
			zap.String("body", string(body)),
		)
		m.next.ServeHTTP(rw, r)
		logger.Info(
			"",
			zap.String("trace_id", traceId),
			zap.String("type", "response"),
			zap.String("method", r.Method),
			zap.String("uri", r.URL.Path),
			zap.String("host", r.Host),
			zap.Object("header", (*helper.ZapHeader)(&r.Header)),
			zap.String("body", string(rw.Body)),
			zap.Int("status_code", rw.StatusCode),
		)
	})
}

func NewLog() Log {
	return Log{}
}
